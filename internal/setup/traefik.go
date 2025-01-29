package setup

import (
	"errors"
	"os"

	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/swarm"
	"github.com/mainanick/shellby/internal/constants"
	"github.com/mainanick/shellby/internal/utils"
	"github.com/traefik/traefik/v3/pkg/config/dynamic"
	"github.com/traefik/traefik/v3/pkg/config/static"
	acmeprovider "github.com/traefik/traefik/v3/pkg/provider/acme"
	"github.com/traefik/traefik/v3/pkg/provider/docker"
	"github.com/traefik/traefik/v3/pkg/provider/file"
	yaml "gopkg.in/yaml.v3"
)

func swarmSpec() swarm.ServiceSpec {
	return swarm.ServiceSpec{
		TaskTemplate: swarm.TaskSpec{
			ContainerSpec: &swarm.ContainerSpec{
				Image: "traefik:v3.3.2",
				Mounts: []mount.Mount{
					{
						Type:   mount.TypeBind,
						Source: "/var/run/docker.sock",
						Target: "/var/run/docker.sock",
					},
					{
						Type:   mount.TypeBind,
						Source: constants.TRAEFIK_FILE,
						Target: "/etc/traefik/traefik.yml",
					},
					{
						Type:   mount.TypeBind,
						Source: constants.DYNAMIC_TRAEFIK_PATH,
						Target: constants.DYNAMIC_TRAEFIK_PATH,
					},
				},
			},
			Networks: []swarm.NetworkAttachmentConfig{
				{
					Target: "shellby-network",
				},
			},
			Placement: &swarm.Placement{
				Constraints: []string{"node.role == manager"},
			},
		},
		Mode: swarm.ServiceMode{
			Replicated: &swarm.ReplicatedService{
				Replicas: utils.Uint64Ptr(1),
			},
		},
		EndpointSpec: &swarm.EndpointSpec{
			Ports: []swarm.PortConfig{
				{
					TargetPort:    80,
					PublishedPort: 80,
					PublishMode:   swarm.PortConfigPublishModeHost,
				},
				{
					TargetPort:    443,
					PublishedPort: 443,
					PublishMode:   swarm.PortConfigPublishModeHost,
				},
				{
					// Dashboard
					TargetPort:    8080,
					PublishedPort: 8080,
					PublishMode:   swarm.PortConfigPublishModeHost,
				},
			},
		},
	}

}

func CreateServerTraefikConfig() error {
	cfg := dynamic.Configuration{
		HTTP: &dynamic.HTTPConfiguration{
			Routers: map[string]*dynamic.Router{
				"shellby-router": {
					EntryPoints: []string{"web"},
					Service:     "shellby-service",
					Rule:        "Host('shellby.docker.localhost') && PathPrefix(`/`)",
				},
			},
			Services: map[string]*dynamic.Service{
				"shellby-service": {
					LoadBalancer: &dynamic.ServersLoadBalancer{
						Servers: []dynamic.Server{
							{
								URL: "http://shellby:3100",
							},
						},
						PassHostHeader: utils.BoolPtr(true),
					},
				},
			},
		},
	}

	yamlBytes, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	if _, err := os.Stat(constants.SHELLBY_TRAEFIK_FILE); errors.Is(err, os.ErrNotExist) {
		return os.WriteFile(constants.SHELLBY_TRAEFIK_FILE, yamlBytes, 0644)
	}
	return nil

}

func CreateTraefikConfig() error {
	cfg := static.Configuration{
		Providers: &static.Providers{
			Docker: &docker.Provider{
				Shared: docker.Shared{
					Watch:            true,
					ExposedByDefault: true,
					Network:          "shellby-network",
				},
			},
			Swarm: &docker.SwarmProvider{
				Shared: docker.Shared{
					ExposedByDefault: true,
					Watch:            true,
				},
			},
			File: &file.Provider{
				Directory: constants.DYNAMIC_TRAEFIK_PATH,
				Watch:     true,
			},
		},
		EntryPoints: map[string]*static.EntryPoint{
			"web": {
				Address: ":80",
			},
			"websecure": {
				Address: ":443",
				HTTP: static.HTTPConfig{
					TLS: &static.TLSConfig{
						CertResolver: "letsencrypt",
					},
				},
			},
		},
		API: &static.API{
			Insecure: true,
		},
		CertificatesResolvers: map[string]static.CertificateResolver{
			"letsencrypt": {
				ACME: &acmeprovider.Configuration{
					Email:   "",
					Storage: constants.DYNAMIC_TRAEFIK_PATH + "/acme.json",
					HTTPChallenge: &acmeprovider.HTTPChallenge{
						EntryPoint: "web",
					},
				},
			},
		},
	}
	yamlBytes, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	if _, err := os.Stat(constants.TRAEFIK_FILE); errors.Is(err, os.ErrNotExist) {
		return os.WriteFile(constants.TRAEFIK_FILE, yamlBytes, 0644)
	}
	return nil
}

func SetupTraefik() {
	err := CreateTraefikConfig()
	if err != nil {
		panic(err)
	}

	err = CreateServerTraefikConfig()
	if err != nil {
		panic(err)
	}
	// docker.Client.ServiceCreate(context.TODO(), swarmSpec(), types.ServiceCreateOptions{
	// 	QueryRegistry: true,
	// })
}
