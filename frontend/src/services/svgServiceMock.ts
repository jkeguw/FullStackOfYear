// SVG处理服务 - 本地模拟版本
import { parseSvg, getSvgViewBox } from './svgService';
import { mockMice } from '@/data/mockMice';

/**
 * 获取本地SVG鼠标列表
 */
export async function getLocalSvgMouseList() {
  return {
    devices: mockMice,
    total: mockMice.length
  };
}

/**
 * 获取本地鼠标SVG数据
 * @param mouseId 鼠标ID
 * @param view 视图类型 (top或side)
 * @returns SVG内容
 */
export async function getLocalMouseSvgData(mouseId: string, view: 'top' | 'side' = 'top'): Promise<string> {
  const mouse = mockMice.find(m => m.id === mouseId);
  if (!mouse || !mouse.svgData) return '';
  
  return view === 'top' ? (mouse.svgData.topView || '') : (mouse.svgData.sideView || '');
}

/**
 * 创建本地重叠SVG
 */
export async function createLocalOverlaySvg(
  deviceIds: string[],
  view: 'top' | 'side',
  opacityValues: number[],
  colors?: string[]
): Promise<string> {
  try {
    if (!deviceIds || deviceIds.length === 0) {
      console.error('无法创建重叠SVG: 没有提供设备ID');
      return '';
    }
    
    // 获取本地SVG数据
    const svgs: string[] = [];
    for (const id of deviceIds) {
      const mouse = mockMice.find(m => m.id === id);
      if (mouse && mouse.svgData) {
        svgs.push(view === 'top' ? (mouse.svgData.topView || '') : (mouse.svgData.sideView || ''));
      }
    }

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
 * 创建本地并排比较SVG
 */
export async function createLocalSideBySideSvg(
  deviceIds: string[],
  view: 'top' | 'side'
): Promise<string> {
  try {
    // 获取本地SVG数据
    const svgs: string[] = [];
    const labels: string[] = [];
    
    for (const id of deviceIds) {
      const mouse = mockMice.find(m => m.id === id);
      if (mouse && mouse.svgData) {
        svgs.push(view === 'top' ? (mouse.svgData.topView || '') : (mouse.svgData.sideView || ''));
        labels.push(`${mouse.brand} ${mouse.name}`);
      }
    }

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

export default {
  getLocalSvgMouseList,
  getLocalMouseSvgData,
  createLocalOverlaySvg,
  createLocalSideBySideSvg
};