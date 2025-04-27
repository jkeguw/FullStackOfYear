// SVG处理服务 - 与后端API集成版本
import { getMouseSVG, compareSVGs, getSVGMouseList } from '@/api/device';

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

    if (svgDoc.documentElement instanceof SVGSVGElement) {
      return svgDoc.documentElement;
    }
    
    // 如果不是SVGSVGElement，创建一个新的并复制内容
    const newSvg = document.createElementNS('http://www.w3.org/2000/svg', 'svg');
    const content = svgDoc.documentElement;
    
    // 复制属性
    Array.from(content.attributes).forEach(attr => {
      newSvg.setAttribute(attr.name, attr.value);
    });
    
    // 复制子节点
    // @ts-ignore - SVG Node manipulation type issues
    while (content.firstChild) {
      // @ts-ignore - SVG Node manipulation type issues
      newSvg.appendChild(content.firstChild);
    }
    
    return newSvg;
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
 * 获取鼠标SVG数据
 * @param mouseId 鼠标ID
 * @param view 视图类型 (top或side)
 * @returns SVG内容
 */
export async function getMouseSvgData(mouseId: string, view: 'top' | 'side' = 'top'): Promise<string> {
  try {
    const response = await getMouseSVG(mouseId, view);
    if (response && response.data && response.data.svgData) {
      return response.data.svgData;
    }
    throw new Error('无法获取SVG数据');
  } catch (error) {
    console.error(`获取鼠标SVG数据失败: ${mouseId} (${view})`, error);
    return '';
  }
}

/**
 * 创建重叠的SVG图像用于比较
 * @param deviceIds 设备ID数组
 * @param view 视图类型 (top或side)
 * @param opacityValues 对应的透明度值数组
 * @param colors 对应的颜色数组 (可选)
 * @returns 合并后的SVG字符串
 */
export async function createOverlaySvg(
  deviceIds: string[],
  view: 'top' | 'side',
  opacityValues: number[],
  colors?: string[]
): Promise<string> {
  try {
    // 调用后端API获取SVG比较数据
    const response = await compareSVGs({
      deviceIds: deviceIds,
      view: view
    });

    if (!response || !response.data || !response.data.devices) {
      throw new Error('无法获取比较数据');
    }

    // 获取SVG数据数组
    const svgs = response.data.devices.map(device => device.svgData);

    // 假如我们需要使用前端处理合并SVG，可以使用下面的代码
    // 但此时我们会更倾向于让后端返回已经处理好的SVG
    // 这里作为备用方案保留

    // 解析第一个SVG以获取基本信息
    if (svgs.length === 0) return '';
    const baseSvg = parseSvg(svgs[0]);
    const viewBox = getSvgViewBox(baseSvg);

    // 创建新的SVG容器
    const svgContainer = document.createElementNS('http://www.w3.org/2000/svg', 'svg');
    svgContainer.setAttribute(
      'viewBox',
      `${viewBox.minX} ${viewBox.minY} ${viewBox.width} ${viewBox.height}`
    );
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
        colorMatrix.setAttribute(
          'values',
          `0 0 0 0 ${parseInt(color.slice(1, 3), 16) / 255} 
                                           0 0 0 0 ${parseInt(color.slice(3, 5), 16) / 255} 
                                           0 0 0 0 ${parseInt(color.slice(5, 7), 16) / 255} 
                                           0 0 0 1 0`
        );

        filter.appendChild(colorMatrix);
        svgContainer.appendChild(filter);
        group.setAttribute('filter', `url(#${filterId})`);
      }

      // 将SVG元素的所有子节点复制到组中
      // @ts-ignore - SVG Node manipulation type issues
      while (svgElement.firstChild) {
        // @ts-ignore - SVG Node manipulation type issues
        group.appendChild(svgElement.firstChild);
      }

      svgContainer.appendChild(group);
    });

    return new XMLSerializer().serializeToString(svgContainer);
  } catch (error) {
    console.error('创建重叠SVG失败:', error);
    return '';
  }
}

/**
 * 创建并排比较的SVG
 * @param deviceIds 设备ID数组
 * @param view 视图类型 (top或side)
 * @returns 合并后的SVG字符串
 */
export async function createSideBySideSvg(
  deviceIds: string[],
  view: 'top' | 'side'
): Promise<string> {
  try {
    // 调用后端API获取SVG比较数据
    const response = await compareSVGs({
      deviceIds: deviceIds,
      view: view
    });

    if (!response || !response.data || !response.data.devices) {
      throw new Error('无法获取比较数据');
    }

    // 获取SVG数据和标签
    const svgs = response.data.devices.map(device => device.svgData);
    const labels = response.data.devices.map(device => `${device.brand} ${device.deviceName}`);

    // 解析所有SVG以确定最大尺寸
    const parsedSvgs = svgs.map((svgContent) => parseSvg(svgContent));
    const viewBoxes = parsedSvgs.map((svg) => getSvgViewBox(svg));

    // 计算总宽度和最大高度
    const padding = 10; // SVG之间的间距
    const totalWidth =
      viewBoxes.reduce((sum, vb) => sum + vb.width, 0) + (viewBoxes.length - 1) * padding;
    const maxHeight = Math.max(...viewBoxes.map((vb) => vb.height));

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
      // @ts-ignore - SVG Node manipulation type issues
      while (svgElement.firstChild) {
        // @ts-ignore - SVG Node manipulation type issues
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
  } catch (error) {
    console.error('创建并排比较SVG失败:', error);
    return '';
  }
}

/**
 * 获取可用于SVG比较的鼠标列表
 * @param type 设备类型
 * @param brand 品牌
 * @returns 可用鼠标列表
 */
export async function getSvgMouseList(type?: string, brand?: string) {
  try {
    const response = await getSVGMouseList({ type, brand });
    return response.data;
  } catch (error) {
    console.error('获取SVG鼠标列表失败:', error);
    return { devices: [], total: 0 };
  }
}

export default {
  parseSvg,
  getSvgViewBox,
  getMouseSvgData,
  createOverlaySvg,
  createSideBySideSvg,
  getSvgMouseList
};