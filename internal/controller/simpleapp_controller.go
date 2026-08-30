package controller

import (
	"context"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	appv1 "tutorial.kubebuilder.io/simpleapp/api/v1"
	"fmt"
)

type SimpleAppReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=apps.tutorial.kubebuilder.io,resources=simpleapps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps.tutorial.kubebuilder.io,resources=simpleapps/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete

func (r *SimpleAppReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	// Esta wea es como una clase
	app := &appv1.SimpleApp{}

	fmt.Println(app)
	fmt.Println("hola el nombre es: ", app.Name)


	// es literalmente una funcion que se le pasa la variable app y la cambia internamente.
	if err := r.Get(ctx, req.NamespacedName, app); err != nil {
		// por si se cae
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	fmt.Println("hola el nombre es: ", app.Name)
	fmt.Println("replicas: ", app.Spec.Replicas)


	// obtiene las replicas.
	desiredReplicas := app.Spec.Replicas
	if desiredReplicas < 0 {
		desiredReplicas = 0
	}

	deploymentName := app.Name + "-deployment"

	fmt.Println("deployment: ", deploymentName)

	deployment := &appsv1.Deployment{}
	//este trae el deployment que es una variable nativa del kubernetes
	err := r.Get(ctx, types.NamespacedName{Name: deploymentName, Namespace: app.Namespace}, deployment)


	if apierrors.IsNotFound(err) {
		// setea un objeto deployment de kubernetes para la creacion mas adelante
		deployment = desiredDeployment(app, deploymentName, desiredReplicas)

		// tiene que ver con el owner reference marca el deployment como hijo del app simpleApp
		if err := controllerutil.SetControllerReference(app, deployment, r.Scheme); err != nil {
			return ctrl.Result{}, err
		}

		// crea el deployment
		if err := r.Create(ctx, deployment); err != nil {
			return ctrl.Result{}, err
		}

		// **** No cacho revisar mas rato. pero creo que vuelve a ejecutar el reconcile
		return ctrl.Result{RequeueAfter: 2 * time.Second}, nil
	}
	if err != nil {
		return ctrl.Result{}, err
	}

	// revisar si hay cambios en los deployment, por ejemplo si cambio la cantidad de replicas.
	changed := false
	// ve si habian cambios en la replicas
	if deployment.Spec.Replicas == nil || *deployment.Spec.Replicas != desiredReplicas {
		deployment.Spec.Replicas = ptr(desiredReplicas)
		changed = true
	}
	//ve si no tiene contenedores y si el nombre de la imagen cambio
	if len(deployment.Spec.Template.Spec.Containers) == 0 || deployment.Spec.Template.Spec.Containers[0].Image != app.Spec.Image {
		deployment.Spec.Template.Spec.Containers = []corev1.Container{{Name: "app", Image: app.Spec.Image}}
		changed = true
	}

	if changed {
		// pos actualiza el deployment
		if err := r.Update(ctx, deployment); err != nil {
			return ctrl.Result{}, err
		}
		// ya entendi el de arriba es para volver a reconcile es como las colas en kafka.
		return ctrl.Result{Requeue: true}, nil
	}

	// hace la revision si estan bien las replicas. Preguntar a Nelson
	newReady := deployment.Status.ReadyReplicas
	newPhase := "Ready"
	if newReady < desiredReplicas {
		newPhase = "Progressing"
	}

	// verifica si las replicas del deployment estan alineadas con el app
	// tambien verifica el Phase (Estudiar mas)
	if app.Status.ReadyReplicas != newReady || app.Status.Phase != newPhase {
		app.Status.ReadyReplicas = newReady
		app.Status.Phase = newPhase
		if err := r.Status().Update(ctx, app); err != nil {
			return ctrl.Result{}, err
		}
	}

	//es para hacer la prueba de eliminacion si el app dice que quiere 0 replicas eliminar todos los deployment
	if desiredReplicas == 0 {
		if err := r.Delete(ctx, deployment); err != nil && !apierrors.IsNotFound(err) {
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	return ctrl.Result{}, nil
}

func desiredDeployment(app *appv1.SimpleApp, name string, replicas int32) *appsv1.Deployment {
	labels := map[string]string{"app": app.Name}
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: app.Namespace},
		Spec: appsv1.DeploymentSpec{
			Replicas: ptr(replicas),
			Selector: &metav1.LabelSelector{MatchLabels: labels},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Labels: labels},
				Spec:       corev1.PodSpec{Containers: []corev1.Container{{Name: "app", Image: app.Spec.Image}}},
			},
		},
	}
}

func ptr(v int32) *int32 { return &v }

func (r *SimpleAppReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appv1.SimpleApp{}).    // observa el app o el recurso personalizado que se hizo en el create api kubebuilder
		Owns(&appsv1.Deployment{}). // observa los deployments de la instancia yaml en samples y que luego el reconcile ejecuto
		Complete(r)
}
