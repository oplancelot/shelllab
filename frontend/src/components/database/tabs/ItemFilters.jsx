import React from 'react'
import { SidebarPanel, SectionHeader, ScrollList } from '../../ui'
import { getQualityColor } from '../../../utils/wow'

const QUALITY_OPTIONS = [
    { value: '', label: 'Any Quality' },
    { value: '0', label: 'Poor (Gray)' },
    { value: '1', label: 'Common (White)' },
    { value: '2', label: 'Uncommon (Green)' },
    { value: '3', label: 'Rare (Blue)' },
    { value: '4', label: 'Epic (Purple)' },
    { value: '5', label: 'Legendary (Orange)' },
]

const BASIC_STAT_OPTIONS = [
    { value: '', label: '- Select -' },
    { value: 'intellect', label: 'Intellect' },
    { value: 'stamina', label: 'Stamina' },
    { value: 'spirit', label: 'Spirit' },
    { value: 'strength', label: 'Strength' },
    { value: 'agility', label: 'Agility' },
]

const OTHER_STAT_OPTIONS = [
    { value: '', label: '- Select -' },
    { value: 'armor', label: 'Armor' },
    { value: 'defense', label: 'Defense' },
    { value: 'dodge', label: 'Dodge' },
    { value: 'parry', label: 'Parry' },
    { value: 'block', label: 'Block' },
    { value: 'hit', label: 'Hit Rating' },
    { value: 'crit', label: 'Crit Rating' },
    { value: 'attack_power', label: 'Attack Power' },
    { value: 'spell_power', label: 'Spell Power' },
    { value: 'fire_res', label: 'Fire Resistance' },
    { value: 'frost_res', label: 'Frost Resistance' },
    { value: 'nature_res', label: 'Nature Resistance' },
    { value: 'shadow_res', label: 'Shadow Resistance' },
    { value: 'arcane_res', label: 'Arcane Resistance' },
]

function FilterSection({ title, children }) {
    return (
        <div className="mb-4 px-2">
            <div className="text-xs font-bold text-wow-gold uppercase mb-1">{title}</div>
            <div className="space-y-2">
                {children}
            </div>
        </div>
    )
}

function RangeInput({ label, minVal, maxVal, onMinChange, onMaxChange }) {
    return (
        <div className="flex flex-col gap-1">
            {label && <span className="text-gray-400 text-xs">{label}</span>}
            <div className="flex items-center gap-2">
                <input
                    type="number"
                    value={minVal}
                    onChange={(e) => onMinChange(e.target.value)}
                    placeholder="0"
                    className="w-full bg-black/40 border border-gray-700 rounded text-xs px-2 py-1 text-white focus:border-wow-gold outline-none"
                    min="0"
                    max="100"
                />
                <span className="text-gray-500">-</span>
                <input
                    type="number"
                    value={maxVal}
                    onChange={(e) => onMaxChange(e.target.value)}
                    placeholder="60"
                    className="w-full bg-black/40 border border-gray-700 rounded text-xs px-2 py-1 text-white focus:border-wow-gold outline-none"
                    min="0"
                    max="100"
                />
            </div>
        </div>
    )
}

function StatRow({ stat, minValue, maxValue, onStatChange, onMinValueChange, onMaxValueChange, options }) {
    return (
        <div className="flex gap-1 items-center">
            <select
                value={stat}
                onChange={(e) => onStatChange(e.target.value)}
                className="flex-1 bg-black/40 border border-gray-700 rounded text-xs px-1 py-1 text-gray-300 focus:border-wow-gold outline-none"
            >
                {options.map((opt, idx) => (
                    <option 
                        key={opt.value} 
                        value={opt.value}
                        style={{ backgroundColor: idx % 2 === 0 ? '#181818' : '#242424', color: '#e0e0e0' }}
                    >
                        {opt.label}
                    </option>
                ))}
            </select>
            <input
                type="number"
                value={minValue}
                onChange={(e) => onMinValueChange(e.target.value)}
                placeholder="Min"
                className="w-14 bg-black/40 border border-gray-700 rounded text-xs px-1 py-1 text-white focus:border-wow-gold outline-none"
                step="0.1"
                min="0"
            />
            <span className="text-gray-500 text-xs">-</span>
            <input
                type="number"
                value={maxValue}
                onChange={(e) => onMaxValueChange(e.target.value)}
                placeholder="Max"
                className="w-14 bg-black/40 border border-gray-700 rounded text-xs px-1 py-1 text-white focus:border-wow-gold outline-none"
                step="0.1"
                min="0"
            />
        </div>
    )
}

export default function ItemFilters({ filters, onChange, onSearch, onReset }) {
    const updateFilter = (key, value) => {
        onChange({ ...filters, [key]: value })
    }

    const updateStat = (index, field, value) => {
        const newStats = [...(filters.stats || [])]
        if (!newStats[index]) newStats[index] = { stat: '', minVal: '', maxVal: '' }
        newStats[index][field] = value
        onChange({ ...filters, stats: newStats })
    }

    const updateOtherStat = (index, field, value) => {
        const newStats = [...(filters.otherStats || [])]
        if (!newStats[index]) newStats[index] = { stat: '', minVal: '', maxVal: '' }
        newStats[index][field] = value
        onChange({ ...filters, otherStats: newStats })
    }
    
    const handleReset = () => {
        if (onReset) onReset()
    }

    return (
        <SidebarPanel>
             <SectionHeader 
                title="Filters" 
                noSearch={true}
            />
            <ScrollList className="p-2 space-y-4">
                {/* Item Level */}
                <FilterSection title="Item Level">
                    <RangeInput
                        minVal={filters.minIlvl || ''}
                        maxVal={filters.maxIlvl || ''}
                        onMinChange={(v) => onChange({...filters, minIlvl: v})}
                        onMaxChange={(v) => onChange({...filters, maxIlvl: v})}
                    />
                </FilterSection>

                {/* Required Level */}
                <FilterSection title="Required Level">
                    <RangeInput
                        minVal={filters.minRl || ''}
                        maxVal={filters.maxRl || ''}
                        onMinChange={(v) => onChange({...filters, minRl: v})}
                        onMaxChange={(v) => onChange({...filters, maxRl: v})}
                    />
                </FilterSection>

                {/* Quality */}
                <FilterSection title="Quality">
                    <div className="flex flex-wrap gap-1">
                        {['Poor', 'Common', 'Uncommon', 'Rare', 'Epic', 'Legendary'].map((q, i) => {
                            const currentQualities = Array.isArray(filters.quality) ? filters.quality : []
                            const isSelected = currentQualities.includes(i)
                            const color = getQualityColor(i)
                            const isHighContrast = ['Rare', 'Epic'].includes(q)
                            
                            return (
                                <button
                                    key={q}
                                    onClick={() => {
                                        const newQualities = isSelected
                                            ? currentQualities.filter(q => q !== i)
                                            : [...currentQualities, i]
                                        onChange({...filters, quality: newQualities})
                                    }}
                                    className={`
                                        px-2 py-1 text-xs rounded border transition-all duration-200 flex-1 min-w-[45%] text-center font-medium
                                        ${isSelected 
                                            ? (isHighContrast ? 'text-white' : 'text-black') 
                                            : 'bg-black/40 border-gray-700 hover:bg-black/60'}
                                    `}
                                    style={{
                                        backgroundColor: isSelected ? color : undefined,
                                        borderColor: isSelected ? color : undefined,
                                        color: isSelected ? undefined : color,
                                        textShadow: isSelected && isHighContrast ? '0 1px 2px rgba(0,0,0,0.8)' : 'none',
                                        boxShadow: isSelected ? `0 0 15px ${color}66` : 'none'
                                    }}
                                >
                                    {q}
                                </button>
                            )
                        })}
                    </div>
                </FilterSection>

                {/* Basic Stats */}
                <FilterSection title="Stats (Min-Max)">
                    {[0, 1, 2].map(i => (
                        <StatRow
                            key={i}
                            stat={filters.stats?.[i]?.stat || ''}
                            minValue={filters.stats?.[i]?.minVal || ''}
                            maxValue={filters.stats?.[i]?.maxVal || ''}
                            onStatChange={(v) => updateStat(i, 'stat', v)}
                            onMinValueChange={(v) => updateStat(i, 'minVal', v)}
                            onMaxValueChange={(v) => updateStat(i, 'maxVal', v)}
                            options={BASIC_STAT_OPTIONS}
                        />
                    ))}
                </FilterSection>

                {/* Other Stats */}
                <FilterSection title="Other Stats">
                    {[0, 1, 2].map(i => (
                        <StatRow
                            key={i}
                            stat={filters.otherStats?.[i]?.stat || ''}
                            minValue={filters.otherStats?.[i]?.minVal || ''}
                            maxValue={filters.otherStats?.[i]?.maxVal || ''}
                            onStatChange={(v) => updateOtherStat(i, 'stat', v)}
                            onMinValueChange={(v) => updateOtherStat(i, 'minVal', v)}
                            onMaxValueChange={(v) => updateOtherStat(i, 'maxVal', v)}
                            options={OTHER_STAT_OPTIONS}
                        />
                    ))}
                </FilterSection>
                
                {/* Reset Button */}
                <div className="pt-2 flex justify-center">
                    <button 
                        onClick={handleReset}
                        className="w-1/2 bg-red-900/30 border border-red-800 text-red-400 hover:bg-red-800/50 hover:text-white text-xs py-2 rounded transition-colors uppercase font-semibold"
                    >
                        Reset Filters
                    </button>
                </div>
            </ScrollList>
        </SidebarPanel>
    )
}
