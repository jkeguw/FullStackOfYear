// 这是一个占位符文件，保留了最基本的接口以兼容MouseShapeVisualization组件
// 原有的手型测量功能已被移除

export function useMeasurement() {
  // 手型尺寸判断函数
  const getHandSizeCategory = (palm: number, length: number, unit: string): string => {
    // 将输入单位转换为毫米进行计算
    const palmInMm = unit === 'mm' ? palm : unit === 'cm' ? palm * 10 : palm * 25.4;

    // 基于手掌宽度进行简单分类
    if (palmInMm >= 95) return 'large';
    if (palmInMm >= 85) return 'medium';
    return 'small';
  };

  const getHandSizeName = (category: string): string => {
    switch (category) {
      case 'large':
        return '大型手';
      case 'medium':
        return '中型手';
      case 'small':
        return '小型手';
      default:
        return '未知';
    }
  };

  // 占位方法，返回模拟数据
  const fetchMeasurements = async () => {
    return {
      measurements: [],
      total: 0
    };
  };

  const fetchMeasurementStats = async () => {
    return null;
  };

  const saveMeasurement = async () => {
    return false;
  };

  const removeMeasurement = async () => {
    return true;
  };

  return {
    getHandSizeCategory,
    getHandSizeName,
    fetchMeasurements,
    fetchMeasurementStats,
    saveMeasurement,
    removeMeasurement
  };
}
