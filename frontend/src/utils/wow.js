export const getQualityColor = (quality) => {
    const colors = {
        0: '#9d9d9d', // Poor
        1: '#ffffff', // Common
        2: '#1eff00', // Uncommon
        3: '#0070dd', // Rare
        4: '#a335ee', // Epic
        5: '#ff8000', // Legendary
        6: '#e6cc80'  // Artifact
    }
    return colors[quality] || '#ffffff'
}

// Get icon path - default to JPG
export const getIconPath = (icon) => {
    if (!icon) return '/local-icons/inv_misc_questionmark.jpg';
    return `/local-icons/${icon.toLowerCase()}.jpg`;
}

// Get PNG variant of icon path
export const getIconPathPng = (icon) => {
    if (!icon) return '/local-icons/inv_misc_questionmark.jpg';
    return `/local-icons/${icon.toLowerCase()}.png`;
}

// Get Zamimg CDN fallback
export const getIconZamimg = (icon) => {
    if (!icon) return 'https://wow.zamimg.com/images/wow/icons/medium/inv_misc_questionmark.jpg';
    return `https://wow.zamimg.com/images/wow/icons/medium/${icon.toLowerCase()}.jpg`;
}

// Icon error handler factory - use in components
// Fallback chain: local JPG -> local PNG -> Zamimg CDN -> hide/error
export const createIconErrorHandler = (iconName, onAllFailed = null) => {
    return (e) => {
        const src = e.target.src;
        
        // Step 1: If local JPG failed, try local PNG
        if (src.includes('/local-icons/') && src.endsWith('.jpg')) {
            e.target.src = getIconPathPng(iconName);
            return;
        }
        
        // Step 2: If local PNG also failed, try Zamimg CDN
        if (src.includes('/local-icons/') && src.endsWith('.png')) {
            e.target.src = getIconZamimg(iconName);
            return;
        }
        
        // Step 3: All fallbacks exhausted
        if (onAllFailed) {
            onAllFailed(e);
        } else {
            e.target.style.display = 'none';
        }
    };
}

export const formatMoney = (money) => {
    if (!money) return { g: 0, s: 0, c: 0 };
    const g = Math.floor(money / 10000);
    const s = Math.floor((money % 10000) / 100);
    const c = money % 100;
    return { g, s, c };
}
