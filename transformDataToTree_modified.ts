transformDataToTree(data, fieldOrder) {
    if (!data || data.length === 0) return [];
    
    // 定义需要求和的字段
    const sumFields = ['lot_count', 'avg_intel_products', 'avg_intel_foundry_non_atm', 'avg_intel_foundry_atm'];
    
    /**
     * 聚合数据项中的字段值
     * @param {Array} items - 数据项数组
     * @param {String} field - 字段名
     * @returns {any} 聚合后的值
     */
    function aggregateFieldValue(items, field) {
        // 如果是需要求和的字段，计算总和
        if (sumFields.includes(field)) {
            return items.reduce((sum, item) => {
                const value = item[field];
                return sum + (typeof value === 'number' ? value : 0);
            }, 0);
        }
        // 否则返回第一个非空值
        return items[0][field];
    }
    
    /**
     * 找到第一个产生分组的字段索引（有多个不同值的字段）
     * @param {Array} items - 数据项
     * @param {Array} fields - 字段数组
     * @returns {Number} 第一个需要分组的字段索引，-1 表示所有字段都相同
     */
    function findFirstDifferentFieldIndex(items, fields) {
        for (let i = 0; i < fields.length; i++) {
            const field = fields[i];
            const uniqueValues = new Set(items.map(item => item[field]));
            if (uniqueValues.size > 1) {
                return i; // 找到第一个有不同值的字段
            }
        }
        return -1; // 所有字段都相同
    }
    
    /**
     * 递归构建树形结构
     * @param {Array} items - 当前层级的数据项
     * @param {Number} startFieldIndex - 开始处理的字段索引
     * @returns {Object|Array} 处理后的数据
     */
    function buildTree(items, startFieldIndex = 0) {
        if (items.length === 0) return [];
        
        const remainingFields = fieldOrder.slice(startFieldIndex);
        
        // 找到第一个需要分组的字段
        const diffIndex = findFirstDifferentFieldIndex(items, remainingFields);
        
        // 如果所有字段都相同，返回聚合后的数据
        if (diffIndex === -1) {
            // 对于所有字段都相同的情况，创建聚合对象
            if (items.length === 1) {
                return items;
            }
            
            // 多条记录需要聚合
            const aggregatedItem: any = {};
            fieldOrder.forEach(field => {
                aggregatedItem[field] = aggregateFieldValue(items, field);
            });
            
            return [aggregatedItem];
        }
        
        // 计算实际的字段索引
        const actualDiffIndex = startFieldIndex + diffIndex;
        
        // nodeName 是最后一个相同字段（即分组字段的前一个）
        const nodeNameIndex = actualDiffIndex === 0 ? 0 : actualDiffIndex - 1;
        const nodeNameField = fieldOrder[nodeNameIndex];
        
        // 需要分组的字段
        const groupByField = fieldOrder[actualDiffIndex];
        
        // 按分组字段分组
        const groupedByField = {};
        items.forEach(item => {
            const key = item[groupByField];
            if (!groupedByField[key]) {
                groupedByField[key] = [];
            }
            groupedByField[key].push(item);
        });
        
        // 构建结果对象
        const result: any = {};
        
        // 添加所有相同的字段（从开始到 nodeNameIndex）
        // 对于相同的字段，使用聚合函数处理
        for (let i = startFieldIndex; i <= nodeNameIndex; i++) {
            const field = fieldOrder[i];
            const value = aggregateFieldValue(items, field);
            result[field] = value === undefined || value === null ? '' : value;
        }
        
        // 设置 nodeName
        result.nodeName = nodeNameField;
        
        // 构建 children
        result.children = [];
        
        Object.keys(groupedByField).forEach(groupValue => {
            const groupItems = groupedByField[groupValue];
            
            // 递归处理每组数据
            const childResults = buildTree(groupItems, actualDiffIndex);
            
            if (Array.isArray(childResults)) {
                childResults.forEach(childItem => {
                    result.children.push(childItem);
                });
            } else {
                result.children.push(childResults);
            }
        });
        
        return result;
    }
    
    const result = buildTree(data, 0);

    
    // 如果结果是数组，直接返回；否则包装成数组
    return Array.isArray(result) ? result : [result];
}
