// SVG处理服务

/**
 * 解析SVG内容 
 * @param svgContent SVG内容
 * @returns 解析后的SVG对象
 */
export function parseSvg(svgContent: string): SVGSVGElement {
  // 确保SVG内容的路径是相对路径，解决可能的路径问题
  if (svgContent && typeof svgContent === 'string') {
    // 尝试修正SVG图像问题：将可能的相对路径转为绝对路径
    svgContent = svgContent.replace(/href="([^h][^t][^t][^p])/g, 'href="/images/$1');
    svgContent = svgContent.replace(/xlink:href="([^h][^t][^t][^p])/g, 'xlink:href="/images/$1');
  }
  if (!svgContent || typeof svgContent !== 'string') {
    // 创建一个空的SVG元素并返回
    const emptySvg = document.createElementNS('http://www.w3.org/2000/svg', 'svg');
    emptySvg.setAttribute('viewBox', '0 0 100 100');
    emptySvg.setAttribute('width', '100');
    emptySvg.setAttribute('height', '100');
    
    const text = document.createElementNS('http://www.w3.org/2000/svg', 'text');
    text.setAttribute('x', '50');
    text.setAttribute('y', '50');
    text.setAttribute('text-anchor', 'middle');
    text.setAttribute('dominant-baseline', 'middle');
    text.setAttribute('fill', '#999');
    text.textContent = 'No SVG Data';
    
    emptySvg.appendChild(text);
    return emptySvg;
  }
  
  try {
    const parser = new DOMParser();
    const svgDoc = parser.parseFromString(svgContent, 'image/svg+xml');
    
    // 检查是否有解析错误
    const parserError = svgDoc.querySelector('parsererror');
    if (parserError) {
      throw new Error('SVG解析错误');
    }
    
    return svgDoc.documentElement as SVGSVGElement;
  } catch (error) {
    console.error('SVG解析失败:', error);
    
    // 返回带有错误信息的SVG
    const errorSvg = document.createElementNS('http://www.w3.org/2000/svg', 'svg');
    errorSvg.setAttribute('viewBox', '0 0 100 100');
    errorSvg.setAttribute('width', '100');
    errorSvg.setAttribute('height', '100');
    
    const text = document.createElementNS('http://www.w3.org/2000/svg', 'text');
    text.setAttribute('x', '50');
    text.setAttribute('y', '50');
    text.setAttribute('text-anchor', 'middle');
    text.setAttribute('dominant-baseline', 'middle');
    text.setAttribute('fill', '#f56c6c');
    text.textContent = 'SVG Error';
    
    errorSvg.appendChild(text);
    return errorSvg;
  }
}

/**
 * 获取SVG的视图框
 * @param svgElement SVG元素
 * @returns 视图框数据
 */
export function getSvgViewBox(svgElement: SVGElement) {
  const viewBox = svgElement.getAttribute('viewBox');
  if (!viewBox) return { minX: 0, minY: 0, width: 0, height: 0 };

  const [minX, minY, width, height] = viewBox.split(' ').map(Number);
  return { minX, minY, width, height };
}

/**
 * 调整SVG大小以匹配物理尺寸
 * @param svgElement SVG元素 
 * @param physicalWidth 物理宽度(mm)
 * @param physicalHeight 物理高度(mm)
 * @returns 调整后的SVG字符串
 */
export function scaleSvgToPhysicalSize(svgElement: SVGElement, physicalWidth: number, physicalHeight: number): string {
  // 复制SVG元素
  const clone = svgElement.cloneNode(true) as SVGElement;
  
  // 获取原始视图框
  const { width, height } = getSvgViewBox(clone);
  
  // 设置宽高属性和transform属性
  clone.setAttribute('width', `${physicalWidth}mm`);
  clone.setAttribute('height', `${physicalHeight}mm`);
  
  // 如果SVG有preserveAspectRatio属性，保留它
  // 否则设置为"none"以允许自由缩放
  if (!clone.hasAttribute('preserveAspectRatio')) {
    clone.setAttribute('preserveAspectRatio', 'none');
  }
  
  return new XMLSerializer().serializeToString(clone);
}

/**
 * 将SVG转换为标准化的格式，以便比较
 * @param svgContent SVG内容
 * @param physicalDimensions 物理尺寸 {width, height}
 * @returns 标准化的SVG字符串
 */
export function standardizeSvg(svgContent: string, physicalDimensions: {width: number, height: number}): string {
  const svgElement = parseSvg(svgContent);
  return scaleSvgToPhysicalSize(svgElement, physicalDimensions.width, physicalDimensions.height);
}

/**
 * 创建重叠的SVG图像用于比较
 * @param svgs SVG内容数组
 * @param opacityValues 对应的透明度值数组
 * @param colors 对应的颜色数组 (可选)
 * @returns 合并后的SVG字符串
 */
export function createOverlaySvg(
  svgs: string[], 
  opacityValues: number[], 
  colors?: string[]
): string {
  // 解析第一个SVG以获取基本信息
  if (svgs.length === 0) return '';
  
  const baseSvg = parseSvg(svgs[0]);
  const viewBox = getSvgViewBox(baseSvg);
  
  // 创建新的SVG容器
  const svgContainer = document.createElementNS('http://www.w3.org/2000/svg', 'svg');
  svgContainer.setAttribute('viewBox', `${viewBox.minX} ${viewBox.minY} ${viewBox.width} ${viewBox.height}`);
  svgContainer.setAttribute('xmlns', 'http://www.w3.org/2000/svg');
  
  // 添加每个SVG作为组
  svgs.forEach((svgContent, index) => {
    const svgElement = parseSvg(svgContent);
    const group = document.createElementNS('http://www.w3.org/2000/svg', 'g');
    
    // 设置组属性
    group.setAttribute('opacity', opacityValues[index].toString());
    
    // 如果提供了颜色，应用颜色变换
    if (colors && colors[index]) {
      const color = colors[index];
      // 创建并应用颜色滤镜
      const filterId = `colorize-${index}`;
      const filter = document.createElementNS('http://www.w3.org/2000/svg', 'filter');
      filter.setAttribute('id', filterId);
      
      // 使用feColorMatrix滤镜适用颜色
      const colorMatrix = document.createElementNS('http://www.w3.org/2000/svg', 'feColorMatrix');
      colorMatrix.setAttribute('type', 'matrix');
      colorMatrix.setAttribute('values', `0 0 0 0 ${parseInt(color.slice(1, 3), 16)/255} 
                                         0 0 0 0 ${parseInt(color.slice(3, 5), 16)/255} 
                                         0 0 0 0 ${parseInt(color.slice(5, 7), 16)/255} 
                                         0 0 0 1 0`);
      
      filter.appendChild(colorMatrix);
      svgContainer.appendChild(filter);
      group.setAttribute('filter', `url(#${filterId})`);
    }
    
    // 将SVG元素的所有子节点复制到组中
    while (svgElement.firstChild) {
      group.appendChild(svgElement.firstChild);
    }
    
    svgContainer.appendChild(group);
  });
  
  return new XMLSerializer().serializeToString(svgContainer);
}

/**
 * 创建并排比较的SVG
 * @param svgs SVG内容数组 
 * @param labels 标签数组
 * @returns 合并后的SVG字符串
 */
export function createSideBySideSvg(svgs: string[], labels: string[]): string {
  if (svgs.length === 0) return '';

  // 解析所有SVG以确定最大尺寸
  const parsedSvgs = svgs.map(svgContent => parseSvg(svgContent));
  const viewBoxes = parsedSvgs.map(svg => getSvgViewBox(svg));
  
  // 计算总宽度和最大高度
  const padding = 10; // SVG之间的间距
  const totalWidth = viewBoxes.reduce((sum, vb) => sum + vb.width, 0) + (viewBoxes.length - 1) * padding;
  const maxHeight = Math.max(...viewBoxes.map(vb => vb.height));
  
  // 创建新的SVG容器
  const svgContainer = document.createElementNS('http://www.w3.org/2000/svg', 'svg');
  svgContainer.setAttribute('viewBox', `0 0 ${totalWidth} ${maxHeight + 30}`); // 额外高度用于标签
  svgContainer.setAttribute('xmlns', 'http://www.w3.org/2000/svg');
  
  // 添加每个SVG
  let currentX = 0;
  parsedSvgs.forEach((svgElement, index) => {
    const vb = viewBoxes[index];
    
    // 创建组用于放置SVG内容
    const group = document.createElementNS('http://www.w3.org/2000/svg', 'g');
    group.setAttribute('transform', `translate(${currentX}, 0)`);
    
    // 将SVG元素的所有子节点复制到组中
    while (svgElement.firstChild) {
      group.appendChild(svgElement.firstChild);
    }
    
    // 添加标签
    if (labels && labels[index]) {
      const text = document.createElementNS('http://www.w3.org/2000/svg', 'text');
      text.setAttribute('x', (vb.width / 2).toString());
      text.setAttribute('y', (maxHeight + 20).toString());
      text.setAttribute('text-anchor', 'middle');
      text.setAttribute('font-family', 'Arial, sans-serif');
      text.setAttribute('font-size', '12');
      text.textContent = labels[index];
      group.appendChild(text);
    }
    
    svgContainer.appendChild(group);
    currentX += vb.width + padding;
  });
  
  return new XMLSerializer().serializeToString(svgContainer);
}

/**
 * 为SVG添加规格尺寸标注
 * @param svgContent SVG内容
 * @param dimensions 尺寸 {width, height, length}
 * @returns 添加标注后的SVG字符串
 */
export function addDimensionMarkers(svgContent: string, dimensions: { width: number, height: number, length: number }): string {
  const svgElement = parseSvg(svgContent);
  const viewBox = getSvgViewBox(svgElement);
  
  // 创建标注组
  const markersGroup = document.createElementNS('http://www.w3.org/2000/svg', 'g');
  markersGroup.setAttribute('class', 'dimension-markers');
  
  // 添加宽度标注
  const widthMarker = document.createElementNS('http://www.w3.org/2000/svg', 'g');
  const widthLine = document.createElementNS('http://www.w3.org/2000/svg', 'line');
  widthLine.setAttribute('x1', '0');
  widthLine.setAttribute('y1', viewBox.height.toString());
  widthLine.setAttribute('x2', viewBox.width.toString());
  widthLine.setAttribute('y2', viewBox.height.toString());
  widthLine.setAttribute('stroke', '#ff0000');
  widthLine.setAttribute('stroke-width', '1');
  widthLine.setAttribute('stroke-dasharray', '5,2');
  
  const widthText = document.createElementNS('http://www.w3.org/2000/svg', 'text');
  widthText.setAttribute('x', (viewBox.width / 2).toString());
  widthText.setAttribute('y', (viewBox.height + 15).toString());
  widthText.setAttribute('text-anchor', 'middle');
  widthText.setAttribute('font-family', 'Arial, sans-serif');
  widthText.setAttribute('font-size', '10');
  widthText.setAttribute('fill', '#ff0000');
  widthText.textContent = `${dimensions.width} mm`;
  
  widthMarker.appendChild(widthLine);
  widthMarker.appendChild(widthText);
  markersGroup.appendChild(widthMarker);
  
  // 添加高度标注
  const heightMarker = document.createElementNS('http://www.w3.org/2000/svg', 'g');
  const heightLine = document.createElementNS('http://www.w3.org/2000/svg', 'line');
  heightLine.setAttribute('x1', '0');
  heightLine.setAttribute('y1', '0');
  heightLine.setAttribute('x2', '0');
  heightLine.setAttribute('y2', viewBox.height.toString());
  heightLine.setAttribute('stroke', '#0000ff');
  heightLine.setAttribute('stroke-width', '1');
  heightLine.setAttribute('stroke-dasharray', '5,2');
  
  const heightText = document.createElementNS('http://www.w3.org/2000/svg', 'text');
  heightText.setAttribute('x', '-15');
  heightText.setAttribute('y', (viewBox.height / 2).toString());
  heightText.setAttribute('text-anchor', 'middle');
  heightText.setAttribute('transform', `rotate(270, -15, ${viewBox.height / 2})`);
  heightText.setAttribute('font-family', 'Arial, sans-serif');
  heightText.setAttribute('font-size', '10');
  heightText.setAttribute('fill', '#0000ff');
  heightText.textContent = `${dimensions.height} mm`;
  
  heightMarker.appendChild(heightLine);
  heightMarker.appendChild(heightText);
  markersGroup.appendChild(heightMarker);
  
  // 添加长度标注 (对于侧视图)
  if (dimensions.length) {
    const lengthMarker = document.createElementNS('http://www.w3.org/2000/svg', 'g');
    const lengthLine = document.createElementNS('http://www.w3.org/2000/svg', 'line');
    lengthLine.setAttribute('x1', '0');
    lengthLine.setAttribute('y1', '0');
    lengthLine.setAttribute('x2', viewBox.width.toString());
    lengthLine.setAttribute('y2', '0');
    lengthLine.setAttribute('stroke', '#00aa00');
    lengthLine.setAttribute('stroke-width', '1');
    lengthLine.setAttribute('stroke-dasharray', '5,2');
    
    const lengthText = document.createElementNS('http://www.w3.org/2000/svg', 'text');
    lengthText.setAttribute('x', (viewBox.width / 2).toString());
    lengthText.setAttribute('y', '-10');
    lengthText.setAttribute('text-anchor', 'middle');
    lengthText.setAttribute('font-family', 'Arial, sans-serif');
    lengthText.setAttribute('font-size', '10');
    lengthText.setAttribute('fill', '#00aa00');
    lengthText.textContent = `${dimensions.length} mm`;
    
    lengthMarker.appendChild(lengthLine);
    lengthMarker.appendChild(lengthText);
    markersGroup.appendChild(lengthMarker);
  }
  
  svgElement.appendChild(markersGroup);
  return new XMLSerializer().serializeToString(svgElement);
}

export default {
  parseSvg,
  getSvgViewBox,
  scaleSvgToPhysicalSize,
  standardizeSvg,
  createOverlaySvg,
  createSideBySideSvg,
  addDimensionMarkers
};