// 图表
import * as echarts from 'echarts/core'
import { BarChart, LineChart } from 'echarts/charts'
import 'echarts/lib/component/markLine';
import {
  TitleComponent,
  TooltipComponent,
  LegendComponent,
  GridComponent,
  DatasetComponent,
  DatasetComponentOption,
  TransformComponent,
  GraphicComponent,
} from 'echarts/components'
// 标签自动布局，全局过渡动画等特性
import { LabelLayout, UniversalTransition } from 'echarts/features'
// 引入 Canvas 渲染器，注意引入 CanvasRenderer 或者 SVGRenderer 是必须的一步
import { CanvasRenderer } from 'echarts/renderers'

// 注册必须的组件
echarts.use([
  BarChart,
  LineChart,
  LegendComponent,
  TitleComponent,
  TooltipComponent,
  GridComponent,
  DatasetComponent,
  TransformComponent,
  // DatasetComponentOption,
  LabelLayout,
  UniversalTransition,
  CanvasRenderer,
  GraphicComponent,
])

export default echarts
