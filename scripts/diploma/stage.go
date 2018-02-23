package diploma

import (
	"github.com/Sovianum/cooling-course-project/postprocessing/templ"
	"github.com/Sovianum/turbocycle/utils/turbine/radial/profilers"
	"github.com/Sovianum/turbocycle/utils/turbine/radial/profiles"
	"github.com/Sovianum/turbocycle/utils/turbine/geom"
	"gonum.org/v1/gonum/mat"
	"github.com/Sovianum/cooling-course-project/core/profiling"
	"github.com/Sovianum/cooling-course-project/postprocessing/dataframes"
	"github.com/Sovianum/turbocycle/impl/turbine/geometry"
	"github.com/Sovianum/turbocycle/impl/turbine/states"
	"github.com/Sovianum/turbocycle/impl/turbine/nodes"
	"github.com/Sovianum/turbocycle/common"
	states2 "github.com/Sovianum/turbocycle/impl/engine/states"
	"math"
)

func saveProfilingTemplate() {
	var inserter = templ.NewDataInserter(
		templatesDir+"/"+profilingTemplate,
		buildDir+"/"+profilingOut,
	)
	if err := inserter.Insert(nil); err != nil {
		panic(err)
	}
}

func saveProfiles(
	profiler profilers.Profiler,
	geomGen geometry.BladingGeometryGenerator,
	hRelArr []float64,
	dataNames [][]string,
	isRotor bool,
) {
	var profileArr = make([]profiles.BladeProfile, len(hRelArr))
	for i, hRel := range hRelArr {
		profileArr[i] = profiles.NewBladeProfileFromProfiler(
			hRel,
			0.01, 0.01,
			0.2, 0.2,
			profiler,
		)
	}

	var installationAngleArr = make([]float64, len(hRelArr))
	var tRelArr = make([]float64, len(hRelArr))
	var tArr = common.LinSpace(0, 1, 200)

	for i, hRel := range hRelArr {
		installationAngleArr[i] = profiler.InstallationAngle(hRel)
		tRelArr[i] = geometry.TRel(hRel, geomGen)
	}

	var coordinatesArr = make([][][][]float64, len(hRelArr))
	for i, bladeProfile := range profileArr {
		coordinatesArr[i] = make([][][]float64, 2)

		if isRotor {
			bladeProfile.Transform(geom.Reflection(0))
		}
		bladeProfile.Transform(geom.Translation(mat.NewVecDense(2, []float64{-1, 0})))
		if !isRotor {
			bladeProfile.Transform(geom.Rotation(installationAngleArr[i] - math.Pi))
		} else {
			bladeProfile.Transform(geom.Rotation(-installationAngleArr[i]))
		}

		coordinatesArr[i][0] = geom.GetCoordinates(tArr, profiles.CircularSegment(bladeProfile))

		bladeProfile.Transform(geom.Translation(mat.NewVecDense(2, []float64{
			tRelArr[i], 0,
		})))
		coordinatesArr[i][1] = geom.GetCoordinates(tArr, profiles.CircularSegment(bladeProfile))
	}

	for i := range hRelArr {
		for j := 0; j != 2; j++ {
			if err := profiling.SaveMatrix(dataDir+"/"+dataNames[i][j], coordinatesArr[i][j]); err != nil {
				panic(err)
			}
		}
	}
}

func saveAngleData(
	profiler profilers.Profiler,
	triangleExtractor func(hRel float64, profiler profilers.Profiler) states.VelocityTriangle,
	filename string,
) {
	var hRelArr = common.LinSpace(0, 1, hPointNum)

	var angleArr = make([][]float64, hPointNum)
	for i, hRel := range hRelArr {
		var triangle = triangleExtractor(hRel, profiler)
		angleArr[i] = make([]float64, 3)

		angleArr[i][0] = hRel
		angleArr[i][1] = triangle.Alpha()
		angleArr[i][2] = triangle.Beta()
	}

	if err := profiling.SaveMatrix(dataDir+"/"+filename, angleArr); err != nil {
		panic(err)
	}
}

func getGasProfilers(stage nodes.TurbineStageNode, rotorProfiler profilers.Profiler) (inletProfiler, outletProfiler profilers.GasProfiler) {
	var pack = stage.GetDataPack()
	var triangleIn = states.NewInletTriangle(pack.U1, pack.C1, pack.Alpha1)
	var triangleOut = stage.VelocityOutput().GetState().(states.VelocityPortState).Triangle

	var tIn = pack.T1
	var tOut = pack.T2

	var pIn = pack.P1
	var pOut = pack.P2

	var reactivity = stage.Reactivity()
	var gas = stage.GasInput().GetState().(states2.GasPortState).Gas

	inletProfiler = profilers.InletGasProfiler(gas, tIn, pIn, reactivity, triangleIn, rotorProfiler)
	outletProfiler = profilers.OutletGasProfiler(gas, tOut, pOut, reactivity, triangleOut, rotorProfiler)
	return
}

func getRotorProfiler(stage nodes.TurbineStageNode) profilers.Profiler {
	var pack = stage.GetDataPack()
	var profiler = profiling.GetInitedRotorProfiler(
		stage.StageGeomGen().RotorGenerator(),
		pack.RotorInletTriangle,
		pack.RotorOutletTriangle,
	)
	return profiler
}

func getStatorProfiler(stage nodes.TurbineStageNode) profilers.Profiler {
	var pack = stage.GetDataPack()
	var profiler = profiling.GetInitedStatorProfiler(
		stage.StageGeomGen().StatorGenerator(),
		stage.VelocityInput().GetState().(states.VelocityPortState).Triangle,
		pack.RotorInletTriangle,
	)
	return profiler
}

func saveStageTemplate(stage nodes.TurbineStageNode) {
	var inserter = templ.NewDataInserter(
		templatesDir+"/"+stageTemplate,
		buildDir+"/"+stageOut,
	)
	var df, err = dataframes.NewStageDF(stage)
	if err != nil {
		panic(err)
	}
	if err := inserter.Insert(df); err != nil {
		panic(err)
	}
}

func solveParticularStage(stage nodes.TurbineStageNode) {
	if err := stage.Process(); err != nil {
		panic(err)
	}
}