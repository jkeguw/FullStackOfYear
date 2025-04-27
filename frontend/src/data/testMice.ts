import type { MouseDevice } from '@/models/MouseModel';

// 测试用鼠标数据
export const testMice: MouseDevice[] = [
  {
    id: 'finalmouse-vxef1pro',
    name: 'VXEF1 Pro',
    brand: 'Finalmouse',
    type: 'mouse',
    price: 799.99,
    dimensions: {
      length: 123,
      width: 69,
      height: 42,
      weight: 49,
      gripWidth: 62
    },
    weight: 49,
    shape: {
      type: 'ergo',
      humpPlacement: 'center',
      frontFlare: 'moderate',
      sideCurvature: 'pronounced',
      handCompatibility: 'right',
      thumbRest: true
    },
    technical: {
      connectivity: ['wireless'],
      sensor: 'Finalsensor 2.0',
      maxDPI: 26000,
      pollingRate: 8000,
      sideButtons: 2,
      weight: 49,
      battery: {
        type: 'lithium-ion',
        capacity: 500,
        life: 70
      }
    },
    recommended: {
      gameTypes: ['FPS', 'MOBA', 'Battle Royale'],
      gripStyles: ['claw', 'palm', 'hybrid'],
      handSizes: ['medium', 'large'],
      dailyUse: true,
      professional: true
    },
    svgData: {
      topView: '/VXEF1PRO-top.svg',
      sideView: '/VXEF1PRO-side.svg'
    },
    description: '轻量化专业电竞鼠标，采用先进无线技术，适合中大型手掌及爪握和掌握握持方式。'
  },
  {
    id: 'logitech-gpw2',
    name: 'G Pro X Superlight 2',
    brand: 'Logitech',
    type: 'mouse',
    price: 699.99,
    dimensions: {
      length: 125,
      width: 65,
      height: 40,
      weight: 60,
      gripWidth: 59
    },
    weight: 60,
    shape: {
      type: 'symmetrical',
      humpPlacement: 'center',
      frontFlare: 'subtle',
      sideCurvature: 'mild',
      handCompatibility: 'ambidextrous'
    },
    technical: {
      connectivity: ['wireless'],
      sensor: 'HERO 2.0',
      maxDPI: 32000,
      pollingRate: 2000,
      sideButtons: 2,
      weight: 60,
      battery: {
        type: 'lithium-ion',
        capacity: 450,
        life: 95
      }
    },
    recommended: {
      gameTypes: ['FPS', 'MOBA', 'Battle Royale'],
      gripStyles: ['claw', 'fingertip', 'hybrid'],
      handSizes: ['small', 'medium', 'large'],
      dailyUse: true,
      professional: true
    },
    svgData: {
      topView: '/gpw2-top.svg',
      sideView: '/gpw2-side.svg'
    },
    description: '罗技最新一代超轻量游戏鼠标，对称设计，适合各种握持方式和手型。'
  },
  {
    id: 'pulsar-hskpro',
    name: 'HSK+ Pro',
    brand: 'Pulsar',
    type: 'mouse',
    price: 349.99,
    dimensions: {
      length: 111,
      width: 66,
      height: 31,
      weight: 45,
      gripWidth: 66
    },
    weight: 45,
    shape: {
      type: 'fingertip',
      humpPlacement: 'none',
      frontFlare: 'minimal',
      sideCurvature: 'minimal',
      handCompatibility: 'ambidextrous'
    },
    technical: {
      connectivity: ['wireless'],
      sensor: 'PAW3395',
      maxDPI: 26000,
      pollingRate: 1000,
      sideButtons: 2,
      weight: 45,
      battery: {
        type: 'lithium-ion',
        capacity: 300,
        life: 60
      }
    },
    recommended: {
      gameTypes: ['FPS', 'Battle Royale'],
      gripStyles: ['fingertip'],
      handSizes: ['medium', 'large'],
      dailyUse: false,
      professional: true
    },
    svgData: {
      topView: '/hskpro-top.svg',
      sideView: '/hskpro-side.svg'
    },
    description: '专为指尖握持设计的超轻量鼠标，极低的高度和重量使其成为FPS游戏玩家的首选。'
  }
];

// 为服务器响应格式化数据
export function getTestMouseData() {
  return {
    code: 0,
    message: 'success',
    data: {
      total: testMice.length,
      page: 1,
      pageSize: 20,
      devices: testMice.map(mouse => ({
        ...mouse,
        type: 'mouse',
        createdAt: new Date().toISOString(),
        updatedAt: new Date().toISOString()
      }))
    }
  };
}

// 模拟API请求
export function mockApiResponse<T>(data: T, delay = 500): Promise<T> {
  return new Promise(resolve => {
    setTimeout(() => {
      resolve(data);
    }, delay);
  });
}