// 比较服务 - 处理鼠标对比逻辑
import type { MouseDevice } from '@/api/device';
import type { MouseComparisonResult } from '@/types/mouse';
import { hardcodedMice } from '@/data/hardcodedMice';

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
  const avgLength =
    mice.reduce((sum, mouse) => sum + (mouse.dimensions?.length || 0), 0) / mice.length;
  const avgWidth =
    mice.reduce((sum, mouse) => sum + (mouse.dimensions?.width || 0), 0) / mice.length;
  const avgHeight =
    mice.reduce((sum, mouse) => sum + (mouse.dimensions?.height || 0), 0) / mice.length;
  const avgWeight =
    mice.reduce((sum, mouse) => sum + (mouse.dimensions?.weight || 0), 0) / mice.length;

  // 计算每个鼠标与平均值的偏差
  const deviations = mice.map((mouse) => {
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
  return Math.max(0, 100 - avgDeviation * 100);
}

/**
 * 计算形状相似度得分
 * @param mice 要比较的鼠标数组
 */
function calculateShapeSimilarity(mice: MouseDevice[]): number {
  if (mice.length < 2) return 100;

  // 基于SVG图像匹配，俯视图匹配俯视图，侧视图匹配侧视图
  const compareSvgImages = (svg1: string, svg2: string): number => {
    if (!svg1 || !svg2) return 0;
    
    // 简化的SVG相似度算法
    // 对于生产环境，应该实现更复杂的图像相似度算法
    // 这里先实现一个基本的长度比对算法作为临时方案
    
    // 去除所有空白和换行符，只保留实际内容
    const cleanSvg1 = svg1.replace(/\s+/g, '');
    const cleanSvg2 = svg2.replace(/\s+/g, '');
    
    // 计算长度相似度
    const maxLength = Math.max(cleanSvg1.length, cleanSvg2.length);
    const minLength = Math.min(cleanSvg1.length, cleanSvg2.length);
    
    if (maxLength === 0) return 0;
    return (minLength / maxLength) * 100;
  };

  // 比较顶视图和侧视图
  let topViewScore = 0;
  let sideViewScore = 0;
  
  // 比较所有鼠标的顶视图
  if (mice.every(mouse => mouse.svgData?.topView)) {
    const baseTopView = mice[0].svgData?.topView;
    const topViewScores = mice.slice(1).map(mouse => 
      compareSvgImages(baseTopView, mouse.svgData?.topView)
    );
    
    // 计算平均顶视图得分
    topViewScore = topViewScores.reduce((sum, score) => sum + score, 0) / 
                   (topViewScores.length || 1);
  }
  
  // 比较所有鼠标的侧视图
  if (mice.every(mouse => mouse.svgData?.sideView)) {
    const baseSideView = mice[0].svgData?.sideView;
    const sideViewScores = mice.slice(1).map(mouse => 
      compareSvgImages(baseSideView, mouse.svgData?.sideView)
    );
    
    // 计算平均侧视图得分
    sideViewScore = sideViewScores.reduce((sum, score) => sum + score, 0) / 
                    (sideViewScores.length || 1);
  }
  
  // 顶视图和侧视图各占50%的权重
  return (topViewScore * 0.5) + (sideViewScore * 0.5);
}

/**
 * 计算技术参数相似度
 * @param mice 鼠标数组
 */
function calculateTechnicalSimilarity(mice: MouseDevice[]): number {
  if (mice.length < 2) return 100;

  // 获取平均DPI和轮询率
  const avgDPI = mice.reduce((sum, mouse) => sum + (mouse.technical?.maxDPI || 0), 0) / mice.length;
  const avgPollingRate =
    mice.reduce((sum, mouse) => sum + (mouse.technical?.pollingRate || 0), 0) / mice.length;

  // 计算每个鼠标与平均值的偏差
  const deviations = mice.map((mouse) => {
    const dpiDev = Math.abs((mouse.technical?.maxDPI || 0) - avgDPI) / (avgDPI || 1);
    const pollingDev =
      Math.abs((mouse.technical?.pollingRate || 0) - avgPollingRate) / (avgPollingRate || 1);

    // 加权平均(DPI差异权重较低，因为它往往有较大的数值范围)
    return dpiDev * 0.3 + pollingDev * 0.7;
  });

  // 平均偏差
  const avgDeviation = deviations.reduce((sum, dev) => sum + dev, 0) / deviations.length;

  // 将偏差转换为相似度得分(0-100)
  return Math.max(0, 100 - avgDeviation * 100);
}

/**
 * 生成鼠标比较结果
 * @param mice 要比较的鼠标
 */
export function generateComparisonResult(mice: MouseDevice[]): MouseComparisonResult {
  // @ts-ignore - Type inconsistency in MouseComparisonResult
  const differences: MouseComparisonResult['differences'] = {};

  // 如果只有一个鼠标或没有鼠标，直接返回
  if (mice.length <= 1) {
    // @ts-ignore - Type inconsistency in return type
    return {
      mice,
      differences: {},
      similarityScore: 100
    };
  }

  // 计算尺寸差异
  differences['dimensions.length'] = {
    property: '长度',
    values: mice.map((m) => m.dimensions?.length || 0),
    differencePercent: calculatePercentageDifference(
      Math.max(...mice.map((m) => m.dimensions?.length || 0)),
      Math.min(...mice.map((m) => m.dimensions?.length || 0))
    )
  };

  differences['dimensions.width'] = {
    property: '宽度',
    values: mice.map((m) => m.dimensions?.width || 0),
    differencePercent: calculatePercentageDifference(
      Math.max(...mice.map((m) => m.dimensions?.width || 0)),
      Math.min(...mice.map((m) => m.dimensions?.width || 0))
    )
  };

  differences['dimensions.height'] = {
    property: '高度',
    values: mice.map((m) => m.dimensions?.height || 0),
    differencePercent: calculatePercentageDifference(
      Math.max(...mice.map((m) => m.dimensions?.height || 0)),
      Math.min(...mice.map((m) => m.dimensions?.height || 0))
    )
  };

  differences['dimensions.weight'] = {
    property: '重量',
    values: mice.map((m) => m.dimensions?.weight || 0),
    differencePercent: calculatePercentageDifference(
      Math.max(...mice.map((m) => m.dimensions?.weight || 0)),
      Math.min(...mice.map((m) => m.dimensions?.weight || 0))
    )
  };

  // 添加握宽差异，如果存在
  if (mice.some((m) => m.dimensions?.gripWidth !== undefined)) {
    differences['dimensions.gripWidth'] = {
      property: '握宽',
      values: mice.map((m) => m.dimensions?.gripWidth || 0),
      differencePercent: calculatePercentageDifference(
        Math.max(...mice.map((m) => m.dimensions?.gripWidth || 0)),
        Math.min(...(mice.map((m) => m.dimensions?.gripWidth || 0).filter((v) => v > 0) || [1]))
      )
    };
  }

  // 计算DPI差异
  differences['technical.maxDPI'] = {
    property: '最大DPI',
    values: mice.map((m) => m.technical?.maxDPI || 0),
    differencePercent: calculatePercentageDifference(
      Math.max(...mice.map((m) => m.technical?.maxDPI || 0)),
      Math.min(...mice.map((m) => m.technical?.maxDPI || 0))
    )
  };

  // 计算轮询率差异
  differences['technical.pollingRate'] = {
    property: '轮询率',
    values: mice.map((m) => m.technical?.pollingRate || 0),
    differencePercent: calculatePercentageDifference(
      Math.max(...mice.map((m) => m.technical?.pollingRate || 0)),
      Math.min(...mice.map((m) => m.technical?.pollingRate || 0))
    )
  };

  // 计算侧键数量差异
  differences['technical.sideButtons'] = {
    property: '侧键数量',
    values: mice.map((m) => m.technical?.sideButtons || 0),
    differencePercent: calculatePercentageDifference(
      Math.max(...mice.map((m) => m.technical?.sideButtons || 0)),
      Math.min(...mice.map((m) => m.technical?.sideButtons || 0))
    )
  };

  // 计算形状类型
  differences['shape.type'] = {
    property: '形状类型',
    values: mice.map((m) => m.shape?.type || '未知'),
    differencePercent: mice.every((m) => m.shape?.type === mice[0]?.shape?.type) ? 0 : 100
  };

  // 计算握持方式
  differences['recommended.gripStyles'] = {
    property: '推荐握持方式',
    values: mice.map((m) => m.recommended?.gripStyles?.join(', ') || '未指定'),
    differencePercent: 0 // 这里不计算差异百分比，只是显示信息
  };

  // 计算总体相似度得分
  const dimensionsScore = calculateDimensionsSimilarity(mice);
  const shapeScore = calculateShapeSimilarity(mice);
  const technicalScore = calculateTechnicalSimilarity(mice);

  // 修改权重，使形状匹配占主导地位
  const similarityScore = Math.round(
    dimensionsScore * 0.2 + shapeScore * 0.7 + technicalScore * 0.1
  );

  // @ts-ignore - Type inconsistency in return type and MouseDevice definition
  return {
    // @ts-ignore - Type inconsistency between different MouseDevice definitions
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
export function findSimilarMice(
  targetMouse: MouseDevice,
  allMice: MouseDevice[],
  limit = 5
): MouseDevice[] {
  // 过滤掉目标鼠标自身
  const otherMice = allMice.filter((mouse) => mouse.id !== targetMouse.id);

  // 计算每个鼠标与目标鼠标的比较结果
  const comparisonResults = otherMice.map((mouse) => {
    // @ts-ignore - Type inconsistency in MouseDevice definition
    const result = generateComparisonResult([targetMouse, mouse]);
    return {
      mouse,
      similarityScore: result.similarityScore
    };
  });

  // 按相似度降序排序
  comparisonResults.sort((a, b) => b.similarityScore - a.similarityScore);

  // 返回指定数量的结果
  return comparisonResults.slice(0, limit).map((result) => result.mouse);
}

// 添加获取硬编码鼠标数据的辅助函数
export function getHardcodedMice(): MouseDevice[] {
  return hardcodedMice;
}

export default {
  generateComparisonResult,
  findSimilarMice,
  getHardcodedMice
};
