package stub

import (
	"context"
	"fmt"

	"github.com/storageos/storageoscluster-operator/pkg/apis/cluster/v1alpha1"
	"github.com/storageos/storageoscluster-operator/pkg/controller"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/record"

	"github.com/operator-framework/operator-sdk/pkg/sdk"
)

// NewHandler returns a new event handler given a recorder and controller.
func NewHandler(eRec record.EventRecorder, c *controller.ClusterController) sdk.Handler {
	return &Handler{eventRecorder: eRec, controller: c}
}

// Handler contains the controller and event broadcast recorder.
type Handler struct {
	eventRecorder record.EventRecorder
	controller    *controller.ClusterController
}

// Handle calls the controller reconcile method based on the event.
func (h *Handler) Handle(ctx context.Context, event sdk.Event) error {
	switch o := event.Object.(type) {
	case *v1alpha1.StorageOSCluster:

		// Ignore the delete event since the garbage collector will clean up all secondary resources for the CR
		// All secondary resources must have the CR set as their OwnerReference for this to be the case
		if event.Deleted {
			// If the current cluster is deleted, reset current cluster.
			if h.controller.IsCurrentCluster(o) {
				h.controller.ResetCurrentCluster()
			}
			return nil
		}

		// Set as the current cluster if there's no current cluster.
		h.controller.SetCurrentClusterIfNone(o)

		// If the event doesn't belongs to the current cluster, do not reconcile.
		// There must be only a single instance of storageos in a cluster.
		if !h.controller.IsCurrentCluster(o) {
			err := fmt.Errorf("can't create more than one storageos cluster")
			h.eventRecorder.Event(o, corev1.EventTypeWarning, "FailedCreation", err.Error())
			return err
		}

		return h.controller.Reconcile(o, h.eventRecorder)
	}

	return nil
}
