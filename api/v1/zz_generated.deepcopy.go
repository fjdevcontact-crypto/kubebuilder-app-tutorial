// Code generated manually for this teaching project. DO NOT EDIT.
package v1

import runtime "k8s.io/apimachinery/pkg/runtime"

func (in *SimpleApp) DeepCopyInto(out *SimpleApp) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
}
func (in *SimpleApp) DeepCopy() *SimpleApp {
	if in == nil {
		return nil
	}
	out := new(SimpleApp)
	in.DeepCopyInto(out)
	return out
}
func (in *SimpleApp) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}
func (in *SimpleAppList) DeepCopyInto(out *SimpleAppList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		out.Items = make([]SimpleApp, len(in.Items))
		for i := range in.Items {
			in.Items[i].DeepCopyInto(&out.Items[i])
		}
	}
}
func (in *SimpleAppList) DeepCopy() *SimpleAppList {
	if in == nil {
		return nil
	}
	out := new(SimpleAppList)
	in.DeepCopyInto(out)
	return out
}
func (in *SimpleAppList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}
func (in *SimpleAppSpec) DeepCopyInto(out *SimpleAppSpec) { *out = *in }
func (in *SimpleAppSpec) DeepCopy() *SimpleAppSpec {
	if in == nil {
		return nil
	}
	out := new(SimpleAppSpec)
	*out = *in
	return out
}
func (in *SimpleAppStatus) DeepCopyInto(out *SimpleAppStatus) { *out = *in }
func (in *SimpleAppStatus) DeepCopy() *SimpleAppStatus {
	if in == nil {
		return nil
	}
	out := new(SimpleAppStatus)
	*out = *in
	return out
}
