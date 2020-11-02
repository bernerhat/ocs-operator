package defaults

import (
	rook "github.com/rook/rook/pkg/apis/rook.io/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	// appLabelSelectorKey is common value for 'Key' field in 'LabelSelectorRequirement'
	appLabelSelectorKey = "app"
	// DefaultNodeAffinity is the NodeAffinity to be used when labelSelector is nil
	DefaultNodeAffinity = &corev1.NodeAffinity{
		RequiredDuringSchedulingIgnoredDuringExecution: getOcsNodeSelector(),
	}
	// DaemonPlacements map contains the default placement configs for the
	// various OCS daemons
	DaemonPlacements = map[string]rook.Placement{
		"all": rook.Placement{
			Tolerations: []corev1.Toleration{
				getOcsToleration(),
			},
		},

		"mon": rook.Placement{
			PodAntiAffinity: &corev1.PodAntiAffinity{
				PreferredDuringSchedulingIgnoredDuringExecution: []corev1.WeightedPodAffinityTerm{
					getWeightedPodAffinityTerm(100, "rook-ceph-mon"),
				},
			},
		},

		"osd": rook.Placement{
			Tolerations: []corev1.Toleration{
				getOcsToleration(),
			},
			PodAntiAffinity: &corev1.PodAntiAffinity{
				PreferredDuringSchedulingIgnoredDuringExecution: []corev1.WeightedPodAffinityTerm{
					getWeightedPodAffinityTerm(100, "rook-ceph-osd"),
				},
			},
		},

		"osd-prepare": rook.Placement{
			Tolerations: []corev1.Toleration{
				getOcsToleration(),
			},
			PodAntiAffinity: &corev1.PodAntiAffinity{
				PreferredDuringSchedulingIgnoredDuringExecution: []corev1.WeightedPodAffinityTerm{
					getWeightedPodAffinityTerm(100, "rook-ceph-osd-prepare"),
				},
			},
		},

		"osd-tsc": rook.Placement{
			Tolerations: []corev1.Toleration{
				getOcsToleration(),
			},
			TopologySpreadConstraints: []corev1.TopologySpreadConstraint{
				getTopologySpreadConstraintsSpec(1, "rook-ceph-osd"),
			},
		},

		"osd-prepare-tsc": rook.Placement{
			Tolerations: []corev1.Toleration{
				getOcsToleration(),
			},
			TopologySpreadConstraints: []corev1.TopologySpreadConstraint{
				getTopologySpreadConstraintsSpec(1, "rook-ceph-osd-prepare", "rook-ceph-osd"),
			},
		},

		"rgw": rook.Placement{
			Tolerations: []corev1.Toleration{
				getOcsToleration(),
			},
			PodAntiAffinity: &corev1.PodAntiAffinity{
				RequiredDuringSchedulingIgnoredDuringExecution: []corev1.PodAffinityTerm{
					getPodAffinityTerm("rook-ceph-rgw"),
				},
				PreferredDuringSchedulingIgnoredDuringExecution: []corev1.WeightedPodAffinityTerm{
					getWeightedPodAffinityTerm(100, "rook-ceph-rgw"),
				},
			},
		},

		"mds": rook.Placement{
			Tolerations: []corev1.Toleration{
				getOcsToleration(),
			},
			PodAntiAffinity: &corev1.PodAntiAffinity{
				RequiredDuringSchedulingIgnoredDuringExecution: []corev1.PodAffinityTerm{
					getPodAffinityTerm("rook-ceph-mds"),
				},
			},
		},

		"noobaa-core": rook.Placement{
			Tolerations: []corev1.Toleration{
				getOcsToleration(),
			},
		},
	}
)

func getTopologySpreadConstraintsSpec(maxSkew int32, selectorValue ...string) corev1.TopologySpreadConstraint {
	topologySpreadConstraints := corev1.TopologySpreadConstraint{
		MaxSkew:           maxSkew,
		TopologyKey:       corev1.LabelHostname,
		WhenUnsatisfiable: "ScheduleAnyway",
		LabelSelector: &metav1.LabelSelector{
			MatchExpressions: []metav1.LabelSelectorRequirement{
				metav1.LabelSelectorRequirement{
					Key:      appLabelSelectorKey,
					Operator: metav1.LabelSelectorOpIn,
					Values:   selectorValue,
				},
			},
		},
	}

	return topologySpreadConstraints
}

func getWeightedPodAffinityTerm(weight int32, selectorValue ...string) corev1.WeightedPodAffinityTerm {
	WeightedPodAffinityTerm := corev1.WeightedPodAffinityTerm{
		Weight: weight,
		PodAffinityTerm: corev1.PodAffinityTerm{
			LabelSelector: &metav1.LabelSelector{
				MatchExpressions: []metav1.LabelSelectorRequirement{
					metav1.LabelSelectorRequirement{
						Key:      appLabelSelectorKey,
						Operator: metav1.LabelSelectorOpIn,
						Values:   selectorValue,
					},
				},
			},
			TopologyKey: corev1.LabelHostname,
		},
	}
	return WeightedPodAffinityTerm
}

func getPodAffinityTerm(selectorValue ...string) corev1.PodAffinityTerm {
	podAffinityTerm := corev1.PodAffinityTerm{
		LabelSelector: &metav1.LabelSelector{
			MatchExpressions: []metav1.LabelSelectorRequirement{
				metav1.LabelSelectorRequirement{
					Key:      appLabelSelectorKey,
					Operator: metav1.LabelSelectorOpIn,
					Values:   selectorValue,
				},
			},
		},
		TopologyKey: corev1.LabelHostname,
	}
	return podAffinityTerm
}

func getOcsToleration() corev1.Toleration {
	toleration := corev1.Toleration{
		Key:      NodeTolerationKey,
		Operator: corev1.TolerationOpEqual,
		Value:    "true",
		Effect:   corev1.TaintEffectNoSchedule,
	}
	return toleration
}

func getOcsNodeSelector() *corev1.NodeSelector {
	nodeSelector := &corev1.NodeSelector{
		NodeSelectorTerms: []corev1.NodeSelectorTerm{
			corev1.NodeSelectorTerm{
				MatchExpressions: []corev1.NodeSelectorRequirement{
					corev1.NodeSelectorRequirement{
						Key:      NodeAffinityKey,
						Operator: corev1.NodeSelectorOpExists,
					},
				},
			},
		},
	}
	return nodeSelector
}
