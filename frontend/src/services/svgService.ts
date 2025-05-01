// SVG processing service - Integrates with backend API for all SVG operations
import { getMouseSVG, compareSVGs, getSVGMouseList } from '@/api/device';

/**
 * Parse SVG content
 * @param svgContent SVG content
 * @returns Parsed SVG object
 */
export function parseSvg(svgContent: string): SVGSVGElement {
  // Ensure SVG content paths are relative, resolving potential path issues
  if (svgContent && typeof svgContent === 'string') {
    // Try to fix SVG image issues: Convert possible relative paths to absolute paths
    svgContent = svgContent.replace(/href="([^h][^t][^t][^p])/g, 'href="/images/$1');
    svgContent = svgContent.replace(/xlink:href="([^h][^t][^t][^p])/g, 'xlink:href="/images/$1');
  }
  if (!svgContent || typeof svgContent !== 'string') {
    // Create an empty SVG element and return it
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

    // Check if there are any parsing errors
    const parserError = svgDoc.querySelector('parsererror');
    if (parserError) {
      throw new Error('SVG parsing error');
    }

    if (svgDoc.documentElement instanceof SVGSVGElement) {
      return svgDoc.documentElement;
    }
    
    // If not an SVGSVGElement, create a new one and copy the content
    const newSvg = document.createElementNS('http://www.w3.org/2000/svg', 'svg');
    const content = svgDoc.documentElement;
    
    // Copy attributes
    Array.from(content.attributes).forEach(attr => {
      newSvg.setAttribute(attr.name, attr.value);
    });
    
    // Copy child nodes
    // @ts-ignore - SVG Node manipulation type issues
    while (content.firstChild) {
      // @ts-ignore - SVG Node manipulation type issues
      newSvg.appendChild(content.firstChild);
    }
    
    return newSvg;
  } catch (error) {
    console.error('SVG parsing failed:', error);

    // Return SVG with error message
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
 * Get SVG viewBox
 * @param svgElement SVG element
 * @returns ViewBox data
 */
export function getSvgViewBox(svgElement: SVGElement) {
  const viewBox = svgElement.getAttribute('viewBox');
  if (!viewBox) return { minX: 0, minY: 0, width: 0, height: 0 };

  const [minX, minY, width, height] = viewBox.split(' ').map(Number);
  return { minX, minY, width, height };
}

/**
 * Get mouse SVG data
 * @param mouseId Mouse ID
 * @param view View type (top or side)
 * @returns SVG content
 */
export async function getMouseSvgData(mouseId: string, view: 'top' | 'side' = 'top'): Promise<string> {
  try {
    const response = await getMouseSVG(mouseId, view);
    if (response && response.data && response.data.svgData) {
      return response.data.svgData;
    }
    throw new Error('无法从API获取SVG数据');
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
    // 基本参数验证
    if (!deviceIds || deviceIds.length === 0) {
      const error = new Error('无法创建重叠SVG: 没有提供设备ID');
      console.error(error);
      throw error;
    }
    
    // 验证视图类型，确保是'top'或'side'
    if (view !== 'top' && view !== 'side') {
      const error = new Error(`无效的视图类型: ${view}. 必须为 'top' 或 'side'`);
      console.error(error);
      view = view as any === 'side' ? 'side' : 'top'; // 默认回退到top
      console.warn(`使用回退视图类型: ${view}`);
    }
    
    // 详细日志记录请求参数
    console.log('SVG比较请求参数:', {
      deviceIds: deviceIds,
      deviceCount: deviceIds.length,
      view: view,
      opacityValues: opacityValues,
      hasColors: colors ? true : false
    });
    
    // 验证至少有2个设备ID
    if (deviceIds.length < 2) {
      console.warn(`设备ID数量不足: 只有 ${deviceIds.length} 个设备. 单设备模式启用.`);
      
      // 优雅处理单设备情况
      if (deviceIds.length === 1) {
        console.log(`正在获取单设备(${deviceIds[0]})的SVG数据`);
        try {
          const singleResponse = await getMouseSVG(deviceIds[0], view);
          console.log('单一设备SVG响应状态:', singleResponse ? '成功' : '失败');
          
          if (singleResponse && singleResponse.data && singleResponse.data.svgData) {
            // 解析单个SVG
            const singleSvg = parseSvg(singleResponse.data.svgData);
            const viewBox = getSvgViewBox(singleSvg);
            
            // 创建新的SVG容器
            const svgContainer = document.createElementNS('http://www.w3.org/2000/svg', 'svg');
            svgContainer.setAttribute(
              'viewBox',
              `${viewBox.minX} ${viewBox.minY} ${viewBox.width} ${viewBox.height}`
            );
            svgContainer.setAttribute('xmlns', 'http://www.w3.org/2000/svg');
            svgContainer.setAttribute('width', '100%');
            svgContainer.setAttribute('height', '100%');
            
            // 添加到容器中
            const group = document.createElementNS('http://www.w3.org/2000/svg', 'g');
            group.setAttribute('opacity', opacityValues[0] ? opacityValues[0].toString() : '1');
            
            if (colors && colors[0]) {
              group.setAttribute('stroke', colors[0]);
              group.setAttribute('stroke-width', '2');
            }
            
            // @ts-ignore - SVG Node manipulation type issues
            while (singleSvg.firstChild) {
              // @ts-ignore - SVG Node manipulation type issues
              group.appendChild(singleSvg.firstChild);
            }
            
            svgContainer.appendChild(group);
            
            // 添加"单设备模式"文本提示
            const text = document.createElementNS('http://www.w3.org/2000/svg', 'text');
            text.setAttribute('x', (viewBox.width / 2).toString());
            text.setAttribute('y', (viewBox.height + 15).toString());
            text.setAttribute('text-anchor', 'middle');
            text.setAttribute('font-family', 'Arial, sans-serif');
            text.setAttribute('font-size', '12');
            text.setAttribute('fill', '#909399');
            text.textContent = '单设备模式 - 请添加其他设备进行比较';
            svgContainer.appendChild(text);
            
            return new XMLSerializer().serializeToString(svgContainer);
          }
        } catch (singleError) {
          console.error('获取单个设备SVG失败:', singleError);
          // 创建并返回错误提示SVG，而不是抛出错误
          return `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 300 150">
            <rect width="300" height="150" fill="#f5f7fa" stroke="#e4e7ed" stroke-width="1"/>
            <text x="150" y="75" text-anchor="middle" fill="#909399" font-family="Arial, sans-serif" font-size="16">单设备模式</text>
            <text x="150" y="100" text-anchor="middle" fill="#f56c6c" font-family="Arial, sans-serif" font-size="12">加载失败: ${singleError.message || '无法获取SVG数据'}</text>
            <text x="150" y="125" text-anchor="middle" fill="#909399" font-family="Arial, sans-serif" font-size="12">请稍后重试或选择其他设备</text>
          </svg>`;
        }
      }
      
      // 如果没有设备ID，返回空白SVG而不是抛出错误
      if (deviceIds.length === 0) {
        return `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 300 150">
          <rect width="300" height="150" fill="#f5f7fa" stroke="#e4e7ed" stroke-width="1"/>
          <text x="150" y="75" text-anchor="middle" fill="#909399" font-family="Arial, sans-serif" font-size="16">未选择设备</text>
          <text x="150" y="100" text-anchor="middle" fill="#909399" font-family="Arial, sans-serif" font-size="12">请至少选择一个设备</text>
        </svg>`;
      }
      
      console.warn('需要至少2个设备ID进行比较操作，将继续API请求流程');
    }
    
    // 发出API请求
    console.log(`正在请求API比较 ${deviceIds.length} 个设备的SVG (${view}视图)`);
    let response;
    try {
      // 构建请求体，确保完全符合后端期望格式
      // 规范化视图类型为字符串字面量
      const viewParam: 'top' | 'side' = view === 'top' ? 'top' : 'side';
      
      // 确保deviceIds是一个有效的数组
      if (!Array.isArray(deviceIds) || deviceIds.length === 0) {
        throw new Error('Invalid deviceIds: must be a non-empty array');
      }
      
      const requestData = {
        deviceIds: deviceIds.slice(0, 3), // 确保不超过3个设备
        view: viewParam
      };
      
      // 详细记录请求体格式，用于调试
      console.log('SVG比较请求体格式:', JSON.stringify(requestData, null, 2));
      
      response = await compareSVGs(requestData);
      console.log('SVG比较API响应状态码:', response?.code);
      console.log('SVG比较API响应数据:', response?.data);
    } catch (apiError) {
      console.error('SVG比较API请求失败:', apiError);
      // 尝试回退到获取单个设备的数据
      throw new Error(`API请求失败: ${apiError.message || '未知错误'}`);
    }

    // 验证API响应
    if (!response) {
      const error = new Error('API返回空响应');
      console.error(error);
      throw error;
    }
    
    if (!response.data) {
      const error = new Error('API返回响应但没有数据字段');
      console.error(error);
      throw error;
    }

    // 验证设备数据
    if (!response.data.devices || response.data.devices.length === 0) {
      const error = new Error('API返回空设备列表');
      console.error(error);
      throw error;
    }

    // 获取SVG数据数组
    const svgs = response.data.devices.map(device => {
      if (!device.svgData) {
        console.warn(`设备 ${device.id || 'unknown'} 没有SVG数据`);
      }
      return device.svgData || '';
    });
    
    if (svgs.length === 0) {
      const error = new Error('API返回的SVG数据为空');
      console.error(error);
      throw error;
    }
    
    console.log(`成功获取 ${svgs.length} 个SVG数据`);
    
    // 检查是否有至少一个有效的SVG
    if (!svgs[0]) {
      console.warn('API返回的第一个SVG数据为空，但会继续处理');
    }

    // 解析第一个SVG以获取基本信息
    const baseSvg = parseSvg(svgs[0]);
    const viewBox = getSvgViewBox(baseSvg);

    // 创建新的SVG容器
    const svgContainer = document.createElementNS('http://www.w3.org/2000/svg', 'svg');
    svgContainer.setAttribute(
      'viewBox',
      `${viewBox.minX} ${viewBox.minY} ${viewBox.width} ${viewBox.height}`
    );
    svgContainer.setAttribute('xmlns', 'http://www.w3.org/2000/svg');
    svgContainer.setAttribute('width', '100%');
    svgContainer.setAttribute('height', '100%');

    // 添加每个SVG作为组
    svgs.forEach((svgContent, index) => {
      const tempDiv = document.createElement('div');
      tempDiv.innerHTML = svgContent;
      const svgElement = tempDiv.querySelector('svg');
      const pathElement = svgElement?.querySelector('path');
      
      if (pathElement) {
        // 如果有path元素，优先使用path创建轮廓
        const group = document.createElementNS('http://www.w3.org/2000/svg', 'g');
        
        // 设置组属性
        group.setAttribute('opacity', opacityValues[index].toString());
        group.setAttribute('fill', 'none'); // 仅显示轮廓
        
        // 应用颜色
        if (colors && colors[index]) {
          const color = colors[index];
          group.setAttribute('stroke', color);
          group.setAttribute('stroke-width', '2');
        } else {
          group.setAttribute('stroke', '#000000');
          group.setAttribute('stroke-width', '2');
        }
        
        // 复制路径元素
        const clonedPath = pathElement.cloneNode(true);
        group.appendChild(clonedPath);
        svgContainer.appendChild(group);
      } else {
        // 如果没有path元素，使用原始方法
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
      }
    });

    return new XMLSerializer().serializeToString(svgContainer);
  } catch (error) {
    console.error('Failed to create overlay SVG:', error);
    
    // Attempt to fetch local SVG data as fallback
    if (deviceIds && deviceIds.length > 0) {
      try {
        console.log('尝试使用本地SVG数据作为回退...');
        
        // Import hardcoded data from frontend
        const hardcodedMice = await import('@/data/hardcodedMice').catch(() => null);
        if (hardcodedMice && hardcodedMice.default) {
          const allMice = hardcodedMice.default;
          
          // Find matching mice SVGs by ID
          const foundMice = deviceIds.map(id => 
            allMice.find(mouse => mouse.id === id || mouse._id === id)
          ).filter(Boolean);
          
          if (foundMice.length > 0) {
            console.log(`从硬编码数据中找到 ${foundMice.length} 个设备`);
            
            // Create a simple SVG for comparison
            const svgContainer = document.createElementNS('http://www.w3.org/2000/svg', 'svg');
            svgContainer.setAttribute('viewBox', '0 0 300 150');
            svgContainer.setAttribute('xmlns', 'http://www.w3.org/2000/svg');
            svgContainer.setAttribute('width', '100%');
            svgContainer.setAttribute('height', '100%');
            
            // Add a background and message
            const bg = document.createElementNS('http://www.w3.org/2000/svg', 'rect');
            bg.setAttribute('width', '300');
            bg.setAttribute('height', '150');
            bg.setAttribute('fill', '#f5f7fa');
            bg.setAttribute('stroke', '#e4e7ed');
            bg.setAttribute('stroke-width', '1');
            svgContainer.appendChild(bg);
            
            // Add devices found message
            const text1 = document.createElementNS('http://www.w3.org/2000/svg', 'text');
            text1.setAttribute('x', '150');
            text1.setAttribute('y', '50');
            text1.setAttribute('text-anchor', 'middle');
            text1.setAttribute('fill', '#409EFF');
            text1.setAttribute('font-family', 'Arial, sans-serif');
            text1.setAttribute('font-size', '14');
            text1.textContent = `用本地数据找到 ${foundMice.length} 个设备`;
            svgContainer.appendChild(text1);
            
            // Add device names
            const deviceText = document.createElementNS('http://www.w3.org/2000/svg', 'text');
            deviceText.setAttribute('x', '150');
            deviceText.setAttribute('y', '75');
            deviceText.setAttribute('text-anchor', 'middle');
            deviceText.setAttribute('fill', '#606266');
            deviceText.setAttribute('font-family', 'Arial, sans-serif');
            deviceText.setAttribute('font-size', '12');
            deviceText.textContent = foundMice.map(m => m.brand + ' ' + m.name).join(', ');
            svgContainer.appendChild(deviceText);
            
            // Add error message
            const text2 = document.createElementNS('http://www.w3.org/2000/svg', 'text');
            text2.setAttribute('x', '150');
            text2.setAttribute('y', '100');
            text2.setAttribute('text-anchor', 'middle');
            text2.setAttribute('fill', '#F56C6C');
            text2.setAttribute('font-family', 'Arial, sans-serif');
            text2.setAttribute('font-size', '12');
            text2.textContent = `API错误: ${error.message || '未知错误'}`;
            svgContainer.appendChild(text2);
            
            // Add retry suggestion
            const text3 = document.createElementNS('http://www.w3.org/2000/svg', 'text');
            text3.setAttribute('x', '150');
            text3.setAttribute('y', '120');
            text3.setAttribute('text-anchor', 'middle');
            text3.setAttribute('fill', '#909399');
            text3.setAttribute('font-family', 'Arial, sans-serif');
            text3.setAttribute('font-size', '11');
            text3.textContent = '请重试或选择其他设备进行比较';
            svgContainer.appendChild(text3);
            
            return new XMLSerializer().serializeToString(svgContainer);
          }
        }
      } catch (fallbackError) {
        console.error('回退到本地数据失败:', fallbackError);
      }
    }
    
    // Return placeholder SVG if everything fails
    return `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 300 150">
      <rect width="300" height="150" fill="#f5f7fa" stroke="#e4e7ed" stroke-width="1"/>
      <text x="150" y="75" text-anchor="middle" fill="#909399" font-family="Arial, sans-serif" font-size="16">SVG比较功能无法加载</text>
      <text x="150" y="100" text-anchor="middle" fill="#909399" font-family="Arial, sans-serif" font-size="12">错误: ${error.message || '未知错误'}</text>
      <text x="150" y="125" text-anchor="middle" fill="#909399" font-family="Arial, sans-serif" font-size="12">请稍后重试</text>
    </svg>`;
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
    // 基本参数验证
    if (!deviceIds || deviceIds.length === 0) {
      const error = new Error('无法创建并排SVG: 没有提供设备ID');
      console.error(error);
      throw error;
    }
    
    // 验证视图类型，确保是'top'或'side'
    if (view !== 'top' && view !== 'side') {
      const error = new Error(`无效的视图类型: ${view}. 必须为 'top' 或 'side'`);
      console.error(error);
      view = view as any === 'side' ? 'side' : 'top'; // 默认回退到top
      console.warn(`使用回退视图类型: ${view}`);
    }
    
    // 详细日志记录请求参数
    console.log('SVG并排比较请求参数:', {
      deviceIds: deviceIds,
      deviceCount: deviceIds.length,
      view: view
    });
    
    // 验证至少有2个设备ID
    if (deviceIds.length < 2) {
      console.warn(`设备ID数量不足: 只有 ${deviceIds.length} 个设备. 单设备模式启用.`);
      
      // 优雅处理单设备情况
      if (deviceIds.length === 1) {
        console.log(`正在获取单设备(${deviceIds[0]})的SVG数据`);
        try {
          const singleResponse = await getMouseSVG(deviceIds[0], view);
          console.log('单一设备SVG响应状态:', singleResponse ? '成功' : '失败');
          
          if (singleResponse && singleResponse.data && singleResponse.data.svgData) {
            // 解析单个SVG
            const singleSvg = parseSvg(singleResponse.data.svgData);
            const viewBox = getSvgViewBox(singleSvg);
            
            // 创建新的SVG容器，为单个设备添加额外的提示空间
            const svgContainer = document.createElementNS('http://www.w3.org/2000/svg', 'svg');
            svgContainer.setAttribute(
              'viewBox',
              `0 0 ${viewBox.width} ${viewBox.height + 30}`
            );
            svgContainer.setAttribute('xmlns', 'http://www.w3.org/2000/svg');
            svgContainer.setAttribute('width', '100%');
            svgContainer.setAttribute('height', '100%');
            
            // 创建组用于放置SVG内容
            const group = document.createElementNS('http://www.w3.org/2000/svg', 'g');
            
            // @ts-ignore - SVG Node manipulation type issues
            while (singleSvg.firstChild) {
              // @ts-ignore - SVG Node manipulation type issues
              group.appendChild(singleSvg.firstChild);
            }
            
            // 添加到容器中
            svgContainer.appendChild(group);
            
            // 添加标签
            const deviceName = singleResponse.data.deviceName || singleResponse.data.name || 'Unknown Device';
            const labelText = `${singleResponse.data.brand || ''} ${deviceName}`;
            
            const label = document.createElementNS('http://www.w3.org/2000/svg', 'text');
            label.setAttribute('x', (viewBox.width / 2).toString());
            label.setAttribute('y', (viewBox.height + 15).toString());
            label.setAttribute('text-anchor', 'middle');
            label.setAttribute('font-family', 'Arial, sans-serif');
            label.setAttribute('font-size', '14');
            label.textContent = labelText;
            svgContainer.appendChild(label);
            
            // 添加"单设备模式"文本提示
            const text = document.createElementNS('http://www.w3.org/2000/svg', 'text');
            text.setAttribute('x', (viewBox.width / 2).toString());
            text.setAttribute('y', (viewBox.height + 28).toString());
            text.setAttribute('text-anchor', 'middle');
            text.setAttribute('font-family', 'Arial, sans-serif');
            text.setAttribute('font-size', '12');
            text.setAttribute('fill', '#909399');
            text.textContent = '单设备模式 - 请添加其他设备进行比较';
            svgContainer.appendChild(text);
            
            return new XMLSerializer().serializeToString(svgContainer);
          }
        } catch (singleError) {
          console.error('获取单个设备SVG失败:', singleError);
          // 创建并返回错误提示SVG，而不是抛出错误
          return `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 300 150">
            <rect width="300" height="150" fill="#f5f7fa" stroke="#e4e7ed" stroke-width="1"/>
            <text x="150" y="75" text-anchor="middle" fill="#909399" font-family="Arial, sans-serif" font-size="16">单设备模式</text>
            <text x="150" y="100" text-anchor="middle" fill="#f56c6c" font-family="Arial, sans-serif" font-size="12">加载失败: ${singleError.message || '无法获取SVG数据'}</text>
            <text x="150" y="125" text-anchor="middle" fill="#909399" font-family="Arial, sans-serif" font-size="12">请稍后重试或选择其他设备</text>
          </svg>`;
        }
      }
      
      // 如果没有设备ID，返回空白SVG而不是抛出错误
      if (deviceIds.length === 0) {
        return `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 300 150">
          <rect width="300" height="150" fill="#f5f7fa" stroke="#e4e7ed" stroke-width="1"/>
          <text x="150" y="75" text-anchor="middle" fill="#909399" font-family="Arial, sans-serif" font-size="16">未选择设备</text>
          <text x="150" y="100" text-anchor="middle" fill="#909399" font-family="Arial, sans-serif" font-size="12">请至少选择一个设备</text>
        </svg>`;
      }
      
      console.warn('需要至少2个设备ID进行并排比较操作，将继续API请求流程');
    }
    
    // 发出API请求
    console.log(`正在请求API并排比较 ${deviceIds.length} 个设备的SVG (${view}视图)`);
    let response;
    try {
      // 规范化视图类型为字符串字面量
      const viewParam: 'top' | 'side' = view === 'top' ? 'top' : 'side';
      
      // 确保deviceIds是一个有效的数组
      if (!Array.isArray(deviceIds) || deviceIds.length === 0) {
        throw new Error('Invalid deviceIds: must be a non-empty array');
      }
      
      const requestData = {
        deviceIds: deviceIds.slice(0, 3), // 确保不超过3个设备
        view: viewParam
      };
      
      // 详细记录请求体格式，用于调试
      console.log('SVG并排比较请求体格式:', JSON.stringify(requestData, null, 2));
      
      response = await compareSVGs(requestData);
      console.log('SVG并排比较API响应状态码:', response?.code);
    } catch (apiError) {
      console.error('SVG并排比较API请求失败:', apiError);
      
      // 尝试回退到获取单个设备的数据并手动组合
      console.log('尝试回退到单独获取各设备SVG并手动组合');
      try {
        const fallbackSvgs = [];
        const fallbackLabels = [];
        
        // 最多处理3个设备
        const limitedIds = deviceIds.slice(0, 3);
        
        // 并行获取所有设备的SVG
        const svgPromises = limitedIds.map(id => getMouseSVG(id, view));
        const svgResponses = await Promise.all(svgPromises);
        
        for (let i = 0; i < svgResponses.length; i++) {
          const resp = svgResponses[i];
          if (resp && resp.data && resp.data.svgData) {
            fallbackSvgs.push(resp.data.svgData);
            fallbackLabels.push(`${resp.data.brand || ''} ${resp.data.deviceName || resp.data.name || 'Unknown'}`);
          }
        }
        
        if (fallbackSvgs.length >= 2) {
          console.log(`回退成功: 获取到 ${fallbackSvgs.length} 个设备的SVG数据`);
          
          // 解析所有SVG
          const parsedSvgs = fallbackSvgs.map(svgContent => parseSvg(svgContent));
          const viewBoxes = parsedSvgs.map(svg => getSvgViewBox(svg));
          
          // 计算总宽度和最大高度
          const padding = 20; // SVG之间的间距
          const totalWidth = viewBoxes.reduce((sum, vb) => sum + vb.width, 0) + (viewBoxes.length - 1) * padding;
          const maxHeight = Math.max(...viewBoxes.map(vb => vb.height));
          
          // 创建新的SVG容器
          const svgContainer = document.createElementNS('http://www.w3.org/2000/svg', 'svg');
          svgContainer.setAttribute('viewBox', `0 0 ${totalWidth} ${maxHeight + 30}`); // 额外高度用于标签
          svgContainer.setAttribute('xmlns', 'http://www.w3.org/2000/svg');
          svgContainer.setAttribute('width', '100%');
          svgContainer.setAttribute('height', '100%');
          
          // 添加每个SVG
          let currentX = 0;
          parsedSvgs.forEach((svg, index) => {
            const vb = viewBoxes[index];
            const group = document.createElementNS('http://www.w3.org/2000/svg', 'g');
            group.setAttribute('transform', `translate(${currentX}, 0)`);
            
            // 复制所有子节点
            // @ts-ignore - SVG Node manipulation type issues
            while (svg.firstChild) {
              // @ts-ignore - SVG Node manipulation type issues
              group.appendChild(svg.firstChild);
            }
            
            // 添加标签
            if (fallbackLabels[index]) {
              const text = document.createElementNS('http://www.w3.org/2000/svg', 'text');
              text.setAttribute('x', (vb.width / 2).toString());
              text.setAttribute('y', (maxHeight + 20).toString());
              text.setAttribute('text-anchor', 'middle');
              text.setAttribute('font-family', 'Arial, sans-serif');
              text.setAttribute('font-size', '14');
              text.textContent = fallbackLabels[index];
              group.appendChild(text);
            }
            
            svgContainer.appendChild(group);
            currentX += vb.width + padding;
          });
          
          return new XMLSerializer().serializeToString(svgContainer);
        }
      } catch (fallbackError) {
        console.error('回退获取SVG数据失败:', fallbackError);
      }
      
      throw new Error(`API请求失败: ${apiError.message || '未知错误'} (回退也失败)`);
    }

    // 验证API响应
    if (!response) {
      const error = new Error('API返回空响应');
      console.error(error);
      throw error;
    }
    
    if (!response.data) {
      const error = new Error('API返回响应但没有数据字段');
      console.error(error);
      throw error;
    }

    // 验证设备数据
    if (!response.data.devices || response.data.devices.length === 0) {
      const error = new Error('API返回空设备列表');
      console.error(error);
      throw error;
    }

    // 获取SVG数据和标签
    const svgs = response.data.devices.map(device => {
      if (!device.svgData) {
        console.warn(`设备 ${device.id || 'unknown'} 没有SVG数据`);
      }
      return device.svgData || '';
    });
    
    if (svgs.length === 0) {
      const error = new Error('API返回的SVG数据为空');
      console.error(error);
      throw error;
    }
    
    console.log(`成功获取 ${svgs.length} 个SVG数据用于并排比较`);
    
    // 检查是否有至少一个有效的SVG
    if (!svgs[0]) {
      console.warn('API返回的第一个SVG数据为空，但会继续处理');
    }
    
    const labels = response.data.devices.map(device => `${device.brand} ${device.deviceName || device.name}`);

    // 解析所有SVG以确定最大尺寸
    const parsedSvgs = svgs.map((svgContent) => parseSvg(svgContent));
    const viewBoxes = parsedSvgs.map((svg) => getSvgViewBox(svg));

    // 计算总宽度和最大高度
    const padding = 20; // SVG之间的间距
    const totalWidth =
      viewBoxes.reduce((sum, vb) => sum + vb.width, 0) + (viewBoxes.length - 1) * padding;
    const maxHeight = Math.max(...viewBoxes.map((vb) => vb.height));

    // 创建新的SVG容器
    const svgContainer = document.createElementNS('http://www.w3.org/2000/svg', 'svg');
    svgContainer.setAttribute('viewBox', `0 0 ${totalWidth} ${maxHeight + 30}`); // 额外高度用于标签
    svgContainer.setAttribute('xmlns', 'http://www.w3.org/2000/svg');
    svgContainer.setAttribute('width', '100%');
    svgContainer.setAttribute('height', '100%');

    // 添加每个SVG
    let currentX = 0;
    svgs.forEach((svgContent, index) => {
      const tempDiv = document.createElement('div');
      tempDiv.innerHTML = svgContent;
      const svgElement = tempDiv.querySelector('svg');
      const pathElement = svgElement?.querySelector('path');
      
      if (pathElement) {
        // 如果有path元素，优先使用path创建轮廓
        const vb = viewBoxes[index];
        const group = document.createElementNS('http://www.w3.org/2000/svg', 'g');
        group.setAttribute('transform', `translate(${currentX}, 0)`);
        
        // 设置轮廓属性
        const pathGroup = document.createElementNS('http://www.w3.org/2000/svg', 'g');
        pathGroup.setAttribute('fill', 'none');
        pathGroup.setAttribute('stroke', '#000000');
        pathGroup.setAttribute('stroke-width', '2');
        
        // 复制路径元素
        const clonedPath = pathElement.cloneNode(true);
        pathGroup.appendChild(clonedPath);
        group.appendChild(pathGroup);
        
        // 添加标签
        if (labels && labels[index]) {
          const text = document.createElementNS('http://www.w3.org/2000/svg', 'text');
          text.setAttribute('x', (vb.width / 2).toString());
          text.setAttribute('y', (maxHeight + 20).toString());
          text.setAttribute('text-anchor', 'middle');
          text.setAttribute('font-family', 'Arial, sans-serif');
          text.setAttribute('font-size', '14');
          text.textContent = labels[index];
          group.appendChild(text);
        }
        
        svgContainer.appendChild(group);
        currentX += vb.width + padding;
      } else {
        // 如果没有path元素，使用原始方法
        const svgElement = parsedSvgs[index];
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
          text.setAttribute('font-size', '14');
          text.textContent = labels[index];
          group.appendChild(text);
        }
  
        svgContainer.appendChild(group);
        currentX += vb.width + padding;
      }
    });

    return new XMLSerializer().serializeToString(svgContainer);
  } catch (error) {
    console.error('创建并排比较SVG失败:', error);
    
    // Attempt to fetch local SVG data as fallback
    if (deviceIds && deviceIds.length > 0) {
      try {
        console.log('尝试使用本地SVG数据作为回退...');
        
        // Import hardcoded data from frontend
        const hardcodedMice = await import('@/data/hardcodedMice').catch(() => null);
        if (hardcodedMice && hardcodedMice.default) {
          const allMice = hardcodedMice.default;
          
          // Find matching mice SVGs by ID
          const foundMice = deviceIds.map(id => 
            allMice.find(mouse => mouse.id === id || mouse._id === id)
          ).filter(Boolean);
          
          if (foundMice.length > 0) {
            console.log(`从硬编码数据中找到 ${foundMice.length} 个设备`);
            
            // Create a simple SVG for side-by-side display
            const svgContainer = document.createElementNS('http://www.w3.org/2000/svg', 'svg');
            svgContainer.setAttribute('viewBox', '0 0 300 150');
            svgContainer.setAttribute('xmlns', 'http://www.w3.org/2000/svg');
            svgContainer.setAttribute('width', '100%');
            svgContainer.setAttribute('height', '100%');
            
            // Add a background
            const bg = document.createElementNS('http://www.w3.org/2000/svg', 'rect');
            bg.setAttribute('width', '300');
            bg.setAttribute('height', '150');
            bg.setAttribute('fill', '#f5f7fa');
            bg.setAttribute('stroke', '#e4e7ed');
            bg.setAttribute('stroke-width', '1');
            svgContainer.appendChild(bg);
            
            // Add fallback message
            const text1 = document.createElementNS('http://www.w3.org/2000/svg', 'text');
            text1.setAttribute('x', '150');
            text1.setAttribute('y', '50');
            text1.setAttribute('text-anchor', 'middle');
            text1.setAttribute('fill', '#409EFF');
            text1.setAttribute('font-family', 'Arial, sans-serif');
            text1.setAttribute('font-size', '14');
            text1.textContent = `用本地数据找到 ${foundMice.length} 个设备`;
            svgContainer.appendChild(text1);
            
            // Add device names
            const deviceText = document.createElementNS('http://www.w3.org/2000/svg', 'text');
            deviceText.setAttribute('x', '150');
            deviceText.setAttribute('y', '75');
            deviceText.setAttribute('text-anchor', 'middle');
            deviceText.setAttribute('fill', '#606266');
            deviceText.setAttribute('font-family', 'Arial, sans-serif');
            deviceText.setAttribute('font-size', '12');
            deviceText.textContent = foundMice.map(m => m.brand + ' ' + m.name).join(', ');
            svgContainer.appendChild(deviceText);
            
            // Add error message
            const text2 = document.createElementNS('http://www.w3.org/2000/svg', 'text');
            text2.setAttribute('x', '150');
            text2.setAttribute('y', '100');
            text2.setAttribute('text-anchor', 'middle');
            text2.setAttribute('fill', '#F56C6C');
            text2.setAttribute('font-family', 'Arial, sans-serif');
            text2.setAttribute('font-size', '12');
            text2.textContent = `API错误: ${error.message || '未知错误'}`;
            svgContainer.appendChild(text2);
            
            // Add retry suggestion
            const text3 = document.createElementNS('http://www.w3.org/2000/svg', 'text');
            text3.setAttribute('x', '150');
            text3.setAttribute('y', '120');
            text3.setAttribute('text-anchor', 'middle');
            text3.setAttribute('fill', '#909399');
            text3.setAttribute('font-family', 'Arial, sans-serif');
            text3.setAttribute('font-size', '11');
            text3.textContent = '无法加载SVG数据，请稍后重试';
            svgContainer.appendChild(text3);
            
            return new XMLSerializer().serializeToString(svgContainer);
          }
        }
      } catch (fallbackError) {
        console.error('回退到本地数据失败:', fallbackError);
      }
    }
    
    // Return placeholder SVG if everything fails
    return `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 300 150">
      <rect width="300" height="150" fill="#f5f7fa" stroke="#e4e7ed" stroke-width="1"/>
      <text x="150" y="75" text-anchor="middle" fill="#909399" font-family="Arial, sans-serif" font-size="16">SVG并排比较功能无法加载</text>
      <text x="150" y="100" text-anchor="middle" fill="#909399" font-family="Arial, sans-serif" font-size="12">错误: ${error.message || '未知错误'}</text>
      <text x="150" y="125" text-anchor="middle" fill="#909399" font-family="Arial, sans-serif" font-size="12">请稍后重试</text>
    </svg>`;
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
    if (response && response.data && response.data.devices) {
      return response.data;
    }
    return { devices: [], total: 0 };
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