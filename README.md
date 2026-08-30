# Kubebuilder SimpleApp — curso mínimo práctico

Este proyecto está pensado para aprender **el patrón general de Kubebuilder** sin distraerse con lógica de negocio compleja.

El caso de uso es deliberadamente simple:

```text
SimpleApp
  spec.replicas
  spec.image
       |
       v
   Reconcile
       |
       v
Deployment
       |
       v
status.readyReplicas
```

Un recurso `SimpleApp` declara una imagen y cantidad de réplicas. El controller mantiene un `Deployment` con ese estado.

## Conceptos que debes aprender aquí

| Concepto | Dónde verlo | Idea mental |
|---|---|---|
| Reconcile | `simpleapp_controller.go` | función que intenta acercar realidad al estado deseado |
| `req ctrl.Request` | inicio de `Reconcile` | trae `namespace/name` del objeto que disparó el evento |
| `r.Get` | dos veces en `Reconcile` | leer objetos desde Kubernetes |
| Spec | `SimpleAppSpec` | lo que el usuario quiere |
| Status | `SimpleAppStatus` | lo que Kubernetes/controller observa |
| `client.Client` | `SimpleAppReconciler` | cliente para Get/Create/Update/Delete |
| `r.Create` | Deployment inexistente | crear recurso hijo |
| `r.Update` | Deployment diferente | modificar recurso hijo |
| `r.Delete` | `replicas: 0` | borrar recurso hijo |
| OwnerReference | `SetControllerReference` | indicar que Deployment pertenece al CR |
| Owns | `SetupWithManager` | cambios del Deployment vuelven a disparar Reconcile |
| SetupWithManager | final controller | conecta controller con el manager |
| Requeue | después de Update | ejecutar Reconcile nuevamente |
| RequeueAfter | después de Create | ejecutar nuevamente después de un tiempo |

## Flujo mental

```text
Evento
  |
  v
Reconcile(req)
  |
  v
r.Get(SimpleApp)
  |
  v
leer Spec
  |
  v
r.Get(Deployment)
  |
  +-- no existe --> SetControllerReference --> r.Create
  |
  +-- existe --> comparar estado deseado/actual
                      |
                      +-- cambió --> r.Update --> Requeue
                      |
                      +-- igual --> actualizar Status
```

## Cómo crear este mismo proyecto con Kubebuilder

Necesitas Go, Docker/Colima, kubectl y kubebuilder.

```bash
mkdir simpleapp
cd simpleapp

kubebuilder init \
  --domain tutorial.kubebuilder.io \
  --repo tutorial.kubebuilder.io/simpleapp

kubebuilder create api \
  --group apps \
  --version v1 \
  --kind SimpleApp
```

Cuando pregunte:

```text
Create Resource [y/n]
y
Create Controller [y/n]
y
```

Después reemplaza el `Spec`, `Status` y controller por los ejemplos de este proyecto.

## Preparar dependencias

Este ZIP ya incluye el CRD y el código mínimo necesario, así que para ejecutarlo basta con:

```bash
go mod tidy
```

En un proyecto generado desde cero con Kubebuilder normalmente también usarías `make generate` y `make manifests`. Aquí los dejamos fuera del camino principal para concentrarnos en entender el controller.

## Instalar el CRD

Primero inicia Kubernetes local, por ejemplo:

```bash
colima start --kubernetes
kubectl get nodes
```

Luego:

```bash
make install
```

Esto aplica `config/crd/simpleapps.yaml` al cluster.

Comprueba:

```bash
kubectl get crd | grep simpleapp
```

## Ejecutar el controller localmente

```bash
make run
```

En este proyecto `make run` equivale a `go run .`.

Déjalo corriendo en una terminal.

## Crear un SimpleApp

En otra terminal:

```bash
kubectl apply -f config/samples/apps_v1_simpleapp.yaml
```

Mira el recurso:

```bash
kubectl get simpleapps
kubectl get simpleapp demo -o yaml
```

Mira el Deployment creado por tu controller:

```bash
kubectl get deployments
kubectl get pods
```

Deberías ver algo parecido a:

```text
demo-deployment   2/2
```

## Probar Update

Edita:

```bash
kubectl edit simpleapp demo
```

Cambia:

```yaml
spec:
  replicas: 4
  image: nginx:1.27
```

Después:

```bash
kubectl get deployment demo-deployment
```

El controller detectará la diferencia y hará `r.Update`.

## Probar cambio de imagen

```bash
kubectl patch simpleapp demo --type merge -p '{"spec":{"image":"nginx:1.28"}}'
```

Luego:

```bash
kubectl get deployment demo-deployment -o jsonpath='{.spec.template.spec.containers[0].image}'
echo
```

## Probar Status

```bash
kubectl get simpleapp demo -o yaml
```

Busca:

```yaml
status:
  readyReplicas: 2
  phase: Ready
```

El `Spec` lo escribe el usuario. El `Status` lo escribe el controller.

## Probar Delete

Este ejemplo usa una regla didáctica:

```text
spec.replicas == 0
    -> r.Delete(Deployment)
```

Ejecuta:

```bash
kubectl patch simpleapp demo --type merge -p '{"spec":{"replicas":0}}'
```

Comprueba:

```bash
kubectl get deployments
```

## Probar OwnerReference

Vuelve a dejar `replicas: 2` y después elimina el CR:

```bash
kubectl delete simpleapp demo
```

Gracias a `SetControllerReference`, Kubernetes sabe que el Deployment pertenece a `SimpleApp` y puede eliminarlo por garbage collection.

```bash
kubectl get deployments
```

## Lo realmente importante que debes memorizar

No memorices cada línea. Memoriza este patrón:

```go
func Reconcile(...) {
    // 1. obtener CR
    r.Get(...)

    // 2. leer Spec

    // 3. buscar recurso hijo
    r.Get(...)

    // 4. si no existe
    r.Create(...)

    // 5. si existe pero está diferente
    r.Update(...)

    // 6. actualizar Status
    r.Status().Update(...)

    // 7. volver a ejecutar cuando corresponda
    Requeue / RequeueAfter
}
```

Y el controller:

```go
func (r *SimpleAppReconciler) SetupWithManager(mgr ctrl.Manager) error {
    return ctrl.NewControllerManagedBy(mgr).
        For(&appv1.SimpleApp{}).
        Owns(&appsv1.Deployment{}).
        Complete(r)
}
```

## Qué estudiar después

Cuando este ejemplo te resulte natural, recién ahí conviene agregar:

1. Conditions en Status.
2. Finalizers.
3. errores y retry/backoff.
4. predicates.
5. múltiples recursos hijos.
6. tests con envtest.
7. webhooks defaulting/validation.

Esos temas son importantes, pero no necesitas aprenderlos antes de entender bien el ciclo básico de `Reconcile`.
