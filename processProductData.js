/**
 * 处理和合并产品数据
 * @param {Array} res - 包含三个对象的数组
 *   res[0]: 包含 total_goal, raw, normal
 *   res[1]: 包含各个产品的 lot 数据（数量）
 *   res[2]: 包含各个产品的 lot 比率数据
 * @returns {Object} 合并后的产品数据对象
 */
function processProductData(res) {
  // 验证输入
  if (!Array.isArray(res) || res.length !== 3) {
    throw new Error('输入必须是包含3个元素的数组');
  }

  const [baseData, lotData, lotRatioData] = res;
  const result = {};

  // 获取所有产品名称（从任何一个数据源中）
  const productNames = new Set([
    ...Object.keys(baseData),
    ...Object.keys(lotData),
    ...Object.keys(lotRatioData)
  ]);

  // 为每个产品构建合并后的数据结构
  productNames.forEach(productName => {
    result[productName] = {
      // 从 res[0] 保留 normal 和 raw
      normal: baseData[productName]?.normal || {},
      raw: baseData[productName]?.raw || {},
      
      // 从 res[1] 获取 lot 数据（数量数据）
      lot: lotData[productName] || {},
      
      // 从 res[2] 获取 lotRaw 数据（比率数据）
      lotRaw: lotRatioData[productName] || {},
      
      // 从 res[0] 获取 total_goal
      total_goal: baseData[productName]?.total_goal || null
    };
  });

  return result;
}

// 示例使用
const testData = [
  {
    "RPLP282": {
      "total_goal": "11.750000",
      "raw": {},
      "normal": {
        "Actual": "5.04",
        "IF_ATM_Down_Time": "4.96",
        "IPG_XRB": "4.96",
        "IF_FSM_DMO_XRB": "4.96",
        "IF_XRB": "4.86",
        "IF_NON_XRB": "4.10",
        "IF_EDISPOSE_REWORK": "2.17"
      }
    },
    "RPRP282": {
      "total_goal": "11.750000",
      "raw": {},
      "normal": {
        "Actual": "4.61",
        "IF_ATM_Down_Time": "4.49",
        "IPG_XRB": "4.49",
        "IF_FSM_DMO_XRB": "4.49",
        "IF_XRB": "4.37",
        "IF_NON_XRB": "3.97",
        "IF_EDISPOSE_REWORK": "2.31"
      }
    },
    "RPLP282/RPRP282": {
      "total_goal": "11.750000",
      "raw": {
        "IF_ATM_Down_Time": "4.81",
        "IPG_XRB": "4.91",
        "IF_FSM_DMO_XRB": "4.90",
        "IF_XRB": "4.78",
        "IF_NON_XRB": "4.27",
        "IF_EDISPOSE_REWORK": "2.37"
      },
      "normal": {
        "Actual": "4.91",
        "IF_ATM_Down_Time": "4.81",
        "IPG_XRB": "4.81",
        "IF_FSM_DMO_XRB": "4.81",
        "IF_XRB": "4.70",
        "IF_NON_XRB": "4.06",
        "IF_EDISPOSE_REWORK": "2.21"
      }
    }
  },
  {
    "RPLP282": {
      "Actual": {
        "Assy": 8253,
        "test_atpo": 8229,
        "test_fpo": 10677,
        "finish": 10913
      },
      "IF_ATM_Down_Time": {
        "Assy": 118,
        "test_atpo": 56,
        "test_fpo": 1,
        "finish": 73
      },
      "IPG_XRB": {
        "Assy": 0,
        "test_atpo": 0,
        "test_fpo": 0,
        "finish": 0
      },
      "IF_FSM_DMO_XRB": {
        "Assy": 8,
        "test_atpo": 2,
        "test_fpo": 1,
        "finish": 8
      },
      "IF_XRB": {
        "Assy": 84,
        "test_atpo": 254,
        "test_fpo": 47,
        "finish": 634
      },
      "IF_NON_XRB": {
        "Assy": 5449,
        "test_atpo": 756,
        "test_fpo": 400,
        "finish": 1415
      },
      "IF_EDISPOSE_REWORK": {
        "test_atpo": 3
      }
    },
    "RPRP282": {
      "Actual": {
        "Assy": 3544,
        "test_atpo": 3328,
        "test_fpo": 4884,
        "finish": 4932
      },
      "IF_ATM_Down_Time": {
        "Assy": 43,
        "test_atpo": 85,
        "test_fpo": 0,
        "finish": 16
      },
      "IPG_XRB": {
        "Assy": 0,
        "test_atpo": 0,
        "test_fpo": 0,
        "finish": 3
      },
      "IF_FSM_DMO_XRB": {
        "Assy": 0,
        "test_atpo": 0,
        "test_fpo": 0,
        "finish": 0
      },
      "IF_XRB": {
        "Assy": 19,
        "test_atpo": 99,
        "test_fpo": 40,
        "finish": 166
      },
      "IF_NON_XRB": {
        "Assy": 2361,
        "test_atpo": 297,
        "test_fpo": 241,
        "finish": 306
      },
      "IF_EDISPOSE_REWORK": {
        "test_atpo": 3
      }
    },
    "RPLP282/RPRP282": {
      "Actual": {
        "Assy": 11797,
        "test_atpo": 11557,
        "test_fpo": 15561,
        "finish": 15845
      },
      "IF_ATM_Down_Time": {
        "Assy": 161,
        "test_atpo": 141,
        "test_fpo": 1,
        "finish": 89
      },
      "IPG_XRB": {
        "Assy": 0,
        "test_atpo": 0,
        "test_fpo": 0,
        "finish": 3
      },
      "IF_FSM_DMO_XRB": {
        "Assy": 8,
        "test_atpo": 2,
        "test_fpo": 1,
        "finish": 8
      },
      "IF_XRB": {
        "Assy": 103,
        "test_atpo": 353,
        "test_fpo": 87,
        "finish": 800
      },
      "IF_NON_XRB": {
        "Assy": 7810,
        "test_atpo": 1053,
        "test_fpo": 641,
        "finish": 1721
      },
      "IF_EDISPOSE_REWORK": {
        "test_atpo": 6
      }
    }
  },
  {
    "RPLP282": {
      "Actual": {
        "Assy": "1.30",
        "test_atpo": "2.31",
        "test_fpo": "0.47",
        "finish": "0.96"
      },
      "IF_ATM_Down_Time": {
        "Assy": "1.28",
        "test_atpo": "2.28",
        "test_fpo": "0.47",
        "finish": "0.93"
      },
      "IPG_XRB": {
        "Assy": "1.30",
        "test_atpo": "2.31",
        "test_fpo": "0.47",
        "finish": "0.96"
      },
      "IF_FSM_DMO_XRB": {
        "Assy": "1.30",
        "test_atpo": "2.31",
        "test_fpo": "0.47",
        "finish": "0.96"
      },
      "IF_XRB": {
        "Assy": "1.30",
        "test_atpo": "2.28",
        "test_fpo": "0.46",
        "finish": "0.88"
      },
      "IF_NON_XRB": {
        "Assy": "0.96",
        "test_atpo": "2.22",
        "test_fpo": "0.42",
        "finish": "0.69"
      },
      "IF_EDISPOSE_REWORK": {
        "test_atpo": "2.31"
      }
    },
    "RPRP282": {
      "Actual": {
        "Assy": "1.05",
        "test_atpo": "2.51",
        "test_fpo": "0.40",
        "finish": "0.65"
      },
      "IF_ATM_Down_Time": {
        "Assy": "1.02",
        "test_atpo": "2.43",
        "test_fpo": "0.40",
        "finish": "0.64"
      },
      "IPG_XRB": {
        "Assy": "1.05",
        "test_atpo": "2.51",
        "test_fpo": "0.40",
        "finish": "0.65"
      },
      "IF_FSM_DMO_XRB": {
        "Assy": "1.05",
        "test_atpo": "2.51",
        "test_fpo": "0.40",
        "finish": "0.65"
      },
      "IF_XRB": {
        "Assy": "1.05",
        "test_atpo": "2.49",
        "test_fpo": "0.40",
        "finish": "0.56"
      },
      "IF_NON_XRB": {
        "Assy": "0.79",
        "test_atpo": "2.42",
        "test_fpo": "0.38",
        "finish": "0.62"
      },
      "IF_EDISPOSE_REWORK": {
        "test_atpo": "2.51"
      }
    },
    "RPLP282/RPRP282": {
      "Actual": {
        "Assy": "1.23",
        "test_atpo": "2.37",
        "test_fpo": "0.45",
        "finish": "0.86"
      },
      "IF_ATM_Down_Time": {
        "Assy": "1.20",
        "test_atpo": "2.32",
        "test_fpo": "0.45",
        "finish": "0.84"
      },
      "IPG_XRB": {
        "Assy": "1.23",
        "test_atpo": "2.37",
        "test_fpo": "0.45",
        "finish": "0.86"
      },
      "IF_FSM_DMO_XRB": {
        "Assy": "1.22",
        "test_atpo": "2.37",
        "test_fpo": "0.45",
        "finish": "0.86"
      },
      "IF_XRB": {
        "Assy": "1.22",
        "test_atpo": "2.34",
        "test_fpo": "0.44",
        "finish": "0.78"
      },
      "IF_NON_XRB": {
        "Assy": "0.91",
        "test_atpo": "2.28",
        "test_fpo": "0.41",
        "finish": "0.67"
      },
      "IF_EDISPOSE_REWORK": {
        "test_atpo": "2.37"
      }
    }
  }
];

// 测试函数
console.log('处理后的数据：');
console.log(JSON.stringify(processProductData(testData), null, 2));

// 导出函数（用于Node.js模块）
if (typeof module !== 'undefined' && module.exports) {
  module.exports = processProductData;
}

// 导出函数（用于ES6模块）
// export default processProductData;
