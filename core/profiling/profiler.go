package profiling

import (
	"github.com/Sovianum/turbocycle/utils/turbine/radial/profilers"
	"github.com/Sovianum/turbocycle/impl/turbine/states"
	"github.com/Sovianum/turbocycle/utils/turbine/radial"
	"github.com/Sovianum/turbocycle/impl/turbine/geometry"
	"github.com/Sovianum/turbocycle/common"
)

func GetInitedStatorProfiler(
	geomGen geometry.BladingGeometryGenerator,
	meanInletTriangle, meanOutletTriangle states.VelocityTriangle,
) profilers.Profiler {
	return profilers.NewProfiler(
		1,
		0.7,

		profilers.NewStatorProfilingBehavior(),
		geomGen,

		meanInletTriangle, meanOutletTriangle,
		radial.NewConstantAbsoluteAngleLaw(),
		radial.NewConstantAbsoluteAngleLaw(),

		func(characteristicAngle, hRel float64) float64 {
			return characteristicAngle
		},
		func(characteristicAngle, hRel float64) float64 {
			return characteristicAngle
		},

		func(hRel float64) float64 {
			return common.InterpTolerate(
				hRel,
				[]float64{0, 1},
				[]float64{common.ToRadians(68), common.ToRadians(34)},
			)
		},

		func(hRel float64) float64 {
			return common.InterpTolerate(
				hRel,
				[]float64{0, 1},
				[]float64{common.ToRadians(30), common.ToRadians(15)},
			)
		},
		func(hRel float64) float64 {
			return common.ToRadians(5)
		},

		func(hRel float64) float64 {
			return 0.7
		},
		func(hRel float64) float64 {
			return 1 / 3
		},
	)
}

func GetInitedRotorProfiler(
	geomGen geometry.BladingGeometryGenerator,
	meanInletTriangle, meanOutletTriangle states.VelocityTriangle,
) profilers.Profiler {
	return profilers.NewProfiler(
		1,
		0.7,

		profilers.NewRotorProfilingBehavior(),
		geomGen,

		meanInletTriangle, meanOutletTriangle,
		radial.NewConstantAbsoluteAngleLaw(),
		radial.NewConstantAbsoluteAngleLaw(),

		func(characteristicAngle, hRel float64) float64 {
			return characteristicAngle
		},
		func(characteristicAngle, hRel float64) float64 {
			return characteristicAngle
		},

		func(hRel float64) float64 {
			return common.InterpTolerate(
				hRel,
				[]float64{0, 1},
				[]float64{common.ToRadians(68), common.ToRadians(34)},
			)
		},

		func(hRel float64) float64 {
			return common.InterpTolerate(
				hRel,
				[]float64{0, 1},
				[]float64{common.ToRadians(30), common.ToRadians(15)},
			)
		},
		func(hRel float64) float64 {
			return common.ToRadians(5)
		},

		func(hRel float64) float64 {
			return 0.7
		},
		func(hRel float64) float64 {
			return 1 / 3
		},
	)
}