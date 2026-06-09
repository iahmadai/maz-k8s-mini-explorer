package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"text/tabwriter"

	"github.com/maz/k8s-mini-explorer/internal/config"
	"github.com/maz/k8s-mini-explorer/internal/k8s"
	"github.com/maz/k8s-mini-explorer/internal/models"
	"github.com/maz/k8s-mini-explorer/internal/server"
	"github.com/maz/k8s-mini-explorer/internal/services"
	"github.com/spf13/cobra"
)

type outputFormat string

const (
	outputJSON  outputFormat = "json"
	outputTable outputFormat = "table"
)

func NewRootCmd() *cobra.Command {
	cfg := config.Load()
	var output string

	root := &cobra.Command{
		Use:   "k8s-explorer",
		Short: "Kubernetes Mini Explorer — inspect pods, deployments, events, and services",
	}

	root.PersistentFlags().StringVar(&cfg.KubeconfigPath, "kubeconfig", cfg.KubeconfigPath, "path to kubeconfig file")
	root.PersistentFlags().StringVar(&cfg.KubeContext, "context", cfg.KubeContext, "kubernetes context")
	root.PersistentFlags().StringVar(&cfg.DefaultNS, "namespace", cfg.DefaultNS, "default kubernetes namespace")
	root.PersistentFlags().BoolVar(&cfg.InCluster, "in-cluster", cfg.InCluster, "use in-cluster config")
	root.PersistentFlags().StringVar(&output, "output", "json", "output format: json or table")

	root.AddCommand(newPodsCmd(cfg, &output))
	root.AddCommand(newDeploymentCmd(cfg, &output))
	root.AddCommand(newEventsCmd(cfg, &output))
	root.AddCommand(newServiceHealthCmd(cfg, &output))
	root.AddCommand(newServeCmd(cfg))

	return root
}

func newPodsCmd(cfg *config.Config, output *string) *cobra.Command {
	return &cobra.Command{
		Use:   "pods",
		Short: "List pods in a namespace",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := k8s.NewClient(cfg)
			if err != nil {
				return err
			}
			resp, err := services.NewPodService(client).List(context.Background(), cfg.DefaultNS)
			if err != nil {
				return err
			}
			return printOutput(*output, resp, printPodsTable)
		},
	}
}

func newDeploymentCmd(cfg *config.Config, output *string) *cobra.Command {
	return &cobra.Command{
		Use:   "deployment [name]",
		Short: "Get deployment status",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := k8s.NewClient(cfg)
			if err != nil {
				return err
			}
			resp, err := services.NewDeploymentService(client).Get(context.Background(), cfg.DefaultNS, args[0])
			if err != nil {
				return err
			}
			return printOutput(*output, resp, nil)
		},
	}
}

func newEventsCmd(cfg *config.Config, output *string) *cobra.Command {
	return &cobra.Command{
		Use:   "events",
		Short: "List events in a namespace",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := k8s.NewClient(cfg)
			if err != nil {
				return err
			}
			resp, err := services.NewEventService(client).List(context.Background(), cfg.DefaultNS)
			if err != nil {
				return err
			}
			return printOutput(*output, resp, printEventsTable)
		},
	}
}

func newServiceHealthCmd(cfg *config.Config, output *string) *cobra.Command {
	return &cobra.Command{
		Use:   "service-health [name]",
		Short: "Check service endpoint health",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := k8s.NewClient(cfg)
			if err != nil {
				return err
			}
			resp, err := services.NewServiceHealthService(client).Get(context.Background(), cfg.DefaultNS, args[0])
			if err != nil {
				return err
			}
			return printOutput(*output, resp, nil)
		},
	}
}

func newServeCmd(cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Start REST API server",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := k8s.NewClient(cfg)
			if err != nil {
				return err
			}
			addr := fmt.Sprintf(":%d", cfg.Port)
			slog.Info("starting server", "addr", addr)
			return http.ListenAndServe(addr, server.NewRouter(client))
		},
	}
	cmd.Flags().IntVar(&cfg.Port, "port", cfg.Port, "HTTP port")
	return cmd
}

func printOutput(format string, data any, tableFn func(any) error) error {
	if outputFormat(format) == outputTable && tableFn != nil {
		return tableFn(data)
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(data)
}

func printPodsTable(data any) error {
	resp, ok := data.(*models.PodListResponse)
	if !ok {
		return fmt.Errorf("invalid data type")
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "NAME\tSTATUS\tREADY\tAGE")
	for _, p := range resp.Pods {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", p.Name, p.Status, p.Ready, p.Age)
	}
	return w.Flush()
}

func printEventsTable(data any) error {
	resp, ok := data.(*models.EventListResponse)
	if !ok {
		return fmt.Errorf("invalid data type")
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "TYPE\tREASON\tOBJECT\tAGE\tMESSAGE")
	for _, e := range resp.Events {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", e.Type, e.Reason, e.InvolvedObject, e.Age, e.Message)
	}
	return w.Flush()
}
