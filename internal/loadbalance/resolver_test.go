package loadbalance_test

import (
	"net"
	"net/url"
	"testing"

	api "github.com/fd1az/proglog/api/v1"
	"github.com/fd1az/proglog/internal/config"
	"github.com/fd1az/proglog/internal/loadbalance"
	"github.com/fd1az/proglog/internal/server"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/attributes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/serviceconfig"
)

func TestResolver(t *testing.T) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)

	tlsConfig, err := config.SetupTLSConfig(config.TLSConfig{
		CertFile:      config.ServerCertFile,
		KeyFile:       config.ServerKeyFile,
		CAFile:        config.CAFile,
		Server:        true,
		ServerAddress: "127.0.0.1",
	})
	require.NoError(t, err)
	serverCreds := credentials.NewTLS(tlsConfig)

	srv, err := server.NewGRPCServer(&server.Config{
		GetServerer: &getServers{}, //<label id="get_servers_mock"/>
	}, grpc.Creds(serverCreds))
	require.NoError(t, err)

	go srv.Serve(l)
	// END: setup_test

	// START: mid_test
	conn := &clientConn{}
	tlsConfig, err = config.SetupTLSConfig(config.TLSConfig{
		CertFile:      config.RootClientCertFile,
		KeyFile:       config.RootClientKeyFile,
		CAFile:        config.CAFile,
		Server:        false,
		ServerAddress: "127.0.0.1",
	})
	require.NoError(t, err)
	clientCreds := credentials.NewTLS(tlsConfig)
	opts := resolver.BuildOptions{
		DialCreds: clientCreds,
	}

	r := &loadbalance.Resolver{}

	address, err := addrToURL(l.Addr())

	require.NoError(t, err)
	_, err = r.Build(
		resolver.Target{
			URL: address,
		},
		conn,
		opts,
	)
	require.NoError(t, err)
	// END: mid_test

	// START: finish_test
	wantState := resolver.State{
		Addresses: []resolver.Address{{
			Addr:       "localhost:9001",
			Attributes: attributes.New("is_leader", true),
		}, {
			Addr:       "localhost:9002",
			Attributes: attributes.New("is_leader", false),
		}},
	}
	require.Equal(t, wantState, conn.state)

	conn.state.Addresses = nil
	r.ResolveNow(resolver.ResolveNowOptions{})
	require.Equal(t, wantState, conn.state)
}

// END: finish_test

// START: mock_deps
type getServers struct{}

func (s *getServers) GetServers() ([]*api.Server, error) {
	return []*api.Server{{
		Id:       "leader",
		RpcAddr:  "localhost:9001",
		IsLeader: true,
	}, {
		Id:      "follower",
		RpcAddr: "localhost:9002",
	}}, nil
}

type clientConn struct {
	resolver.ClientConn
	state resolver.State
}

func (c *clientConn) UpdateState(state resolver.State) error {
	c.state = state
	return nil
}

func (c *clientConn) ReportError(err error) {}

func (c *clientConn) NewAddress(addrs []resolver.Address) {}

func (c *clientConn) NewServiceConfig(config string) {}

func (c *clientConn) ParseServiceConfig(
	config string,
) *serviceconfig.ParseResult {
	return nil
}

func addrToURL(addr net.Addr) (url.URL, error) {
	uri, err := url.Parse("tcp://" + addr.String()) // Asegúrate de que addr.String() devuelva una dirección correcta
	if err != nil {
		return url.URL{}, err
	}

	return *uri, nil
}
