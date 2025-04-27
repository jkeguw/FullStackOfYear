// 比较服务 - 处理鼠标对比逻辑
import type { MouseDevice, MouseComparisonResult } from '@/models/MouseModel';

/**
 * 计算两个数值间的百分比差异
 * @param a 值A
 * @param b 值B
 */
function calculatePercentageDifference(a: number, b: number): number {
  if (a === 0 && b === 0) return 0;
  if (a === 0 || b === 0) return 100;
  
  const max = Math.max(a, b);
  const min = Math.min(a, b);
  
  // 防止除以零
  if (min === 0) return 100;
  
  return ((max - min) / min) * 100;
}

/**
 * 计算鼠标物理尺寸的相似度得分
 * @param mice 要比较的鼠标数组
 */
function calculateDimensionsSimilarity(mice: MouseDevice[]): number {
  if (mice.length < 2) return 100;
  
  // 获取所有维度的平均值
  const avgLength = mice.reduce((sum, mouse) => sum + (mouse.dimensions?.length || 0), 0) / mice.length;
  const avgWidth = mice.reduce((sum, mouse) => sum + (mouse.dimensions?.width || 0), 0) / mice.length;
  const avgHeight = mice.reduce((sum, mouse) => sum + (mouse.dimensions?.height || 0), 0) / mice.length;
  const avgWeight = mice.reduce((sum, mouse) => sum + (mouse.dimensions?.weight || 0), 0) / mice.length;
  
  // 计算每个鼠标与平均值的偏差
  const deviations = mice.map(mouse => {
    const lengthDev = Math.abs((mouse.dimensions?.length || 0) - avgLength) / (avgLength || 1);
    const widthDev = Math.abs((mouse.dimensions?.width || 0) - avgWidth) / (avgWidth || 1);
    const heightDev = Math.abs((mouse.dimensions?.height || 0) - avgHeight) / (avgHeight || 1);
    const weightDev = Math.abs((mouse.dimensions?.weight || 0) - avgWeight) / (avgWeight || 1);
    
    // 返回总偏差
    return (lengthDev + widthDev + heightDev + weightDev) / 4;
  });
  
  // 平均偏差
  const avgDeviation = deviations.reduce((sum, dev) => sum + dev, 0) / deviations.length;
  
  // 将偏差转换为相似度得分(0-100)
  return Math.max(0, 100 - (avgDeviation * 100));
}

/**
 * 计算形状相似度得分
 * @param mice 要比较的鼠标数组
 */
function calculateShapeSimilarity(mice: MouseDevice[]): number {
  if (mice.length < 2) return 100;
  
  // 计算鼠标形状特征的相似性 (简化版)
  // 检查形状类型、隆起位置等是否一致
  
  const totalMatches = mice.reduce((matches, mouse, i) => {
    if (i === 0) return matches; // 跳过第一个鼠标(基准)
    
    const baseShape = mice[0]?.shape || {};
    const currentShape = mouse?.shape || {};
    
    let score = 0;
    // 形状类型匹配
    if (baseShape.type === currentShape.type) score += 1;
    // 隆起位置匹配
    if (baseShape.humpPlacement === currentShape.humpPlacement) score += 1;
    // 前部曲线匹配
    if (baseShape.frontFlare === currentShape.frontFlare) score += 1;
    // 侧面曲率匹配
    if (baseShape.sideCurvature === currentShape.sideCurvature) score += 1;
    // 手型兼容性匹配
    if (baseShape.handCompatibility === currentShape.handCompatibility) score += 1;
    
    // 一个鼠标最多5分(所有属性都匹配)
    return matches + (score / 5);
  }, 0);
  
  // 计算总分数 (满分是每个鼠标都与第一个完全匹配)
  return (totalMatches / (mice.length - 1)) * 100;
}

/**
 * 计算技术参数相似度
 * @param mice 鼠标数组
 */
function calculateTechnicalSimilarity(mice: MouseDevice[]): number {
  if (mice.length < 2) return 100;
  
  // 获取平均DPI和轮询率
  const avgDPI = mice.reduce((sum, mouse) => sum + (mouse.technical?.maxDPI || 0), 0) / mice.length;
  const avgPollingRate = mice.reduce((sum, mouse) => sum + (mouse.technical?.pollingRate || 0), 0) / mice.length;
  
  // 计算每个鼠标与平均值的偏差
  const deviations = mice.map(mouse => {
    const dpiDev = Math.abs((mouse.technical?.maxDPI || 0) - avgDPI) / (avgDPI || 1);
    const pollingDev = Math.abs((mouse.technical?.pollingRate || 0) - avgPollingRate) / (avgPollingRate || 1);
    
    // 加权平均(DPI差异权重较低，因为它往往有较大的数值范围)
    return (dpiDev * 0.3) + (pollingDev * 0.7);
  });
  
  // 平均偏差
  const avgDeviation = deviations.reduce((sum, dev) => sum + dev, 0) / deviations.length;
  
  // 将偏差转换为相似度得分(0-100)
  return Math.max(0, 100 - (avgDeviation * 100));
}

/**
 * 生成鼠标比较结果
 * @param mice 要比较的鼠标
 */
export function generateComparisonResult(mice: MouseDevice[]): MouseComparisonResult {
  const differences: MouseComparisonResult['differences'] = {};
  
  // 如果只有一个鼠标或没有鼠标，直接返回
  if (mice.length <= 1) {
    return {
      mice,
      differences: {},
      similarityScore: 100
    };
  }
  
  // 计算尺寸差异
  differences['dimensions.length'] = {
    property: '长度',
    values: mice.map(m => m.dimensions?.length || 0),
    differencePercent: calculatePercentageDifference(
      Math.max(...mice.map(m => m.dimensions?.length || 0)),
      Math.min(...mice.map(m => m.dimensions?.length || 0))
    )
  };
  
  differences['dimensions.width'] = {
    property: '宽度',
    values: mice.map(m => m.dimensions?.width || 0),
    differencePercent: calculatePercentageDifference(
      Math.max(...mice.map(m => m.dimensions?.width || 0)),
      Math.min(...mice.map(m => m.dimensions?.width || 0))
    )
  };
  
  differences['dimensions.height'] = {
    property: '高度',
    values: mice.map(m => m.dimensions?.height || 0),
    differencePercent: calculatePercentageDifference(
      Math.max(...mice.map(m => m.dimensions?.height || 0)),
      Math.min(...mice.map(m => m.dimensions?.height || 0))
    )
  };
  
  differences['dimensions.weight'] = {
    property: '重量',
    values: mice.map(m => m.dimensions?.weight || 0),
    differencePercent: calculatePercentageDifference(
      Math.max(...mice.map(m => m.dimensions?.weight || 0)),
      Math.min(...mice.map(m => m.dimensions?.weight || 0))
    )
  };
  
  // 添加握宽差异，如果存在
  if (mice.some(m => m.dimensions?.gripWidth !== undefined)) {
    differences['dimensions.gripWidth'] = {
      property: '握宽',
      values: mice.map(m => m.dimensions?.gripWidth || 0),
      differencePercent: calculatePercentageDifference(
        Math.max(...mice.map(m => m.dimensions?.gripWidth || 0)),
        Math.min(...mice.map(m => m.dimensions?.gripWidth || 0).filter(v => v > 0) || [1])
      )
    };
  }
  
  // 计算DPI差异
  differences['technical.maxDPI'] = {
    property: '最大DPI',
    values: mice.map(m => m.technical?.maxDPI || 0),
    differencePercent: calculatePercentageDifference(
      Math.max(...mice.map(m => m.technical?.maxDPI || 0)),
      Math.min(...mice.map(m => m.technical?.maxDPI || 0))
    )
  };
  
  // 计算轮询率差异
  differences['technical.pollingRate'] = {
    property: '轮询率',
    values: mice.map(m => m.technical?.pollingRate || 0),
    differencePercent: calculatePercentageDifference(
      Math.max(...mice.map(m => m.technical?.pollingRate || 0)),
      Math.min(...mice.map(m => m.technical?.pollingRate || 0))
    )
  };
  
  // 计算侧键数量差异
  differences['technical.sideButtons'] = {
    property: '侧键数量',
    values: mice.map(m => m.technical?.sideButtons || 0),
    differencePercent: calculatePercentageDifference(
      Math.max(...mice.map(m => m.technical?.sideButtons || 0)),
      Math.min(...mice.map(m => m.technical?.sideButtons || 0))
    )
  };
  
  // 计算形状类型
  differences['shape.type'] = {
    property: '形状类型',
    values: mice.map(m => m.shape?.type || '未知'),
    differencePercent: mice.every(m => m.shape?.type === mice[0]?.shape?.type) ? 0 : 100
  };
  
  // 计算握持方式
  differences['recommended.gripStyles'] = {
    property: '推荐握持方式',
    values: mice.map(m => m.recommended?.gripStyles?.join(', ') || '未指定'),
    differencePercent: 0 // 这里不计算差异百分比，只是显示信息
  };
  
  // 计算总体相似度得分
  const dimensionsScore = calculateDimensionsSimilarity(mice);
  const shapeScore = calculateShapeSimilarity(mice);
  const technicalScore = calculateTechnicalSimilarity(mice);
  
  // 加权计算总相似度
  const similarityScore = Math.round(
    (dimensionsScore * 0.4) + (shapeScore * 0.4) + (technicalScore * 0.2)
  );
  
  return {
    mice,
    differences,
    similarityScore
  };
}

/**
 * 寻找与目标鼠标相似的鼠标
 * @param targetMouse 目标鼠标
 * @param allMice 所有鼠标数组
 * @param limit 限制结果数量
 */
export function findSimilarMice(targetMouse: MouseDevice, allMice: MouseDevice[], limit = 5): MouseDevice[] {
  // 过滤掉目标鼠标自身
  const otherMice = allMice.filter(mouse => mouse.id !== targetMouse.id);
  
  // 计算每个鼠标与目标鼠标的比较结果
  const comparisonResults = otherMice.map(mouse => {
    const result = generateComparisonResult([targetMouse, mouse]);
    return {
      mouse,
      similarityScore: result.similarityScore
    };
  });
  
  // 按相似度降序排序
  comparisonResults.sort((a, b) => b.similarityScore - a.similarityScore);
  
  // 返回指定数量的结果
  return comparisonResults.slice(0, limit).map(result => result.mouse);
}

export default {
  generateComparisonResult,
  findSimilarMice
};