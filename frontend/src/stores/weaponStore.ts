import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

// WeaponCode interface matching the updated Go struct
export interface WeaponCode {
  id: string
  mode: '烽火地带' | '全面战场'
  name: string
  tier: string
  price: number | null      // 改装价格（万）
  build: string             // 改装描述
  code: string              // 改枪码
  range: number | null      // 有效射程（米）
  update_time: string | null // 更新时间
  source: string            // 数据来源: "刀仔" or "武器大师"
}

// Helper function to infer weapon type from weapon name
const getWeaponType = (name: string): string => {
  const typeMap: Record<string, string> = {
    'M4A1': '突击步枪', 'MK47': '突击步枪', 'K416': '突击步枪', 'KC17': '突击步枪',
    'K437': '突击步枪', 'M4': '突击步枪', 'As-Val': '突击步枪', 'ASh-12': '突击步枪',
    'SCAR-H': '突击步枪', 'AK-12': '突击步枪', 'AK-47': '突击步枪', 'FAMAS': '突击步枪',
    'AUG': '突击步枪', 'QBZ': '突击步枪', 'QBZ-95': '突击步枪', 'Type-20': '突击步枪',
    'MP5': '冲锋枪', 'MP7': '冲锋枪', 'MPX': '冲锋枪', 'P90': '冲锋枪',
    'Vector': '冲锋枪', 'UZI': '冲锋枪', 'MAC-10': '冲锋枪', 'Skorpion': '冲锋枪',
    'M14': '射手步枪', 'Mk14': '射手步枪', 'SR-25': '射手步枪', 'G28': '射手步枪',
    'SCAR-HSSR': '射手步枪', 'DMR': '射手步枪', 'SVD': '射手步枪',
    'M1911': '手枪', 'Glock': '手枪', 'P226': '手枪', 'DesertEagle': '手枪',
    'Rex': '手枪', 'Magnum': '手枪', 'M9': '手枪', '93R': '手枪',
    'AWM': '狙击', 'M200': '狙击', 'M24': '狙击', 'Kar98k': '狙击',
    'Mosin': '狙击', 'Lee-Enfield': '狙击', 'Lynx': '狙击', 'TAC-50': '狙击',
    'Marlin': '狙击',
    'M250': '机枪', 'M249': '机枪', 'PKM': '机枪', 'MG42': '机枪',
    'M870': '霰弹枪', 'S12K': '霰弹枪', 'DBS': '霰弹枪', 'Shorty': '霰弹枪',
    'Origin-12': '霰弹枪', 'AA-12': '霰弹枪',
    'Crossbow': '弓弩', 'CompoundBow': '弓弩', 'Arbalist': '弓弩',
  }

  const nameUpper = name.toUpperCase()

  // Check exact matches first
  for (const [key, type] of Object.entries(typeMap)) {
    if (nameUpper.includes(key.toUpperCase())) {
      return type
    }
  }

  // Fallback patterns
  if (nameUpper.includes('冲锋枪') || nameUpper.includes('SMG')) return '冲锋枪'
  if (nameUpper.includes('射手步枪') || nameUpper.includes('DMR')) return '射手步枪'
  if (nameUpper.includes('手枪') || nameUpper.includes('PISTOL')) return '手枪'
  if (nameUpper.includes('狙击') || nameUpper.includes('SNIPER')) return '狙击'
  if (nameUpper.includes('机枪') || nameUpper.includes('MACHINEGUN') || nameUpper.includes('LMG')) return '机枪'
  if (nameUpper.includes('霰弹枪') || nameUpper.includes('SHOTGUN')) return '霰弹枪'
  if (nameUpper.includes('弓') || nameUpper.includes('CROSSBOW') || nameUpper.includes('ARBALIST')) return '弓弩'
  if (nameUpper.includes('步枪') || nameUpper.includes('RIFLE')) return '突击步枪'

  return '其他'
}

export const useWeaponStore = defineStore('weapons', () => {
  // State
  const codes = ref<WeaponCode[]>([])
  const loading = ref(false)
  const searchQuery = ref('')
  const selectedMode = ref<string>('烽火地带') // '烽火地带', '全面战场'
  const selectedWeaponType = ref<string>('all') // 'all', '突击步枪', '冲锋枪', etc.
  const selectedDataSource = ref<string>('刀仔') // '刀仔', '武器大师'

  // Load weapon codes from backend
  const loadCodes = async () => {
    // Check if Wails runtime is ready
    if (!(window as any).go || !(window as any).go.main || !(window as any).go.main.App) {
      console.error('Wails runtime not ready')
      codes.value = []
      return
    }

    loading.value = true
    try {
      let result: WeaponCode[]

      // Load based on selected data source
      if (selectedDataSource.value === '刀仔') {
        result = await (window as any).go.main.App.GetWeaponCodesFromDaoZai()
      } else {
        result = await (window as any).go.main.App.GetWeaponCodesFromWeaponMaster()
      }

      // Ensure result is always an array
      codes.value = Array.isArray(result) ? result : []
    } catch (error) {
      console.error('Failed to load weapon codes:', error)
      codes.value = []
    } finally {
      loading.value = false
    }
  }

  // Get all unique modes
  const uniqueModes = computed(() => {
    const modes = new Set<string>()
    const codesArray = Array.isArray(codes.value) ? codes.value : []
    codesArray.forEach(code => {
      modes.add(code.mode)
    })
    return Array.from(modes)
  })

  // Computed property for filtered codes based on search and filters
  const filteredCodes = computed(() => {
    // Ensure codes.value is always an array
    let result = Array.isArray(codes.value) ? codes.value : []

    // Filter by mode
    result = result.filter(code => code.mode === selectedMode.value)

    // Filter by weapon type
    if (selectedWeaponType.value !== 'all') {
      result = result.filter(code => {
        const weaponType = getWeaponType(code.name)
        return weaponType === selectedWeaponType.value
      })
    }

    // Search query
    if (searchQuery.value) {
      const query = searchQuery.value.toLowerCase()
      result = result.filter(code =>
        code.name.toLowerCase().includes(query) ||
        code.build.toLowerCase().includes(query) ||
        code.code.toLowerCase().includes(query) ||
        code.tier.toLowerCase().includes(query)
      )
    }

    return result
  })

  // Tier priority order
  const tierPriority: Record<string, number> = {
    'T0': 0,
    'T1': 1,
    'T2': 2,
  }

  // Group codes by weapon name for better display
  const groupedCodes = computed(() => {
    const groups = new Map<string, WeaponCode[]>()
    filteredCodes.value.forEach(code => {
      if (!groups.has(code.name)) {
        groups.set(code.name, [])
      }
      groups.get(code.name)!.push(code)
    })

    // Sort by tier (T0 > T1 > T2), then by weapon name (Chinese pinyin)
    // Items without tier (null, '-', '') go to the end
    return Array.from(groups.entries())
      .sort(([nameA, codesA], [nameB, codesB]) => {
        // Get tier from first code of each weapon
        const tierA = codesA[0]?.tier || ''
        const tierB = codesB[0]?.tier || ''

        // Helper to check if tier is valid (not empty, not '-', not 'null')
        const isValidTier = (tier: string) => tier && tier !== '-' && tier !== 'null'

        const hasTierA = isValidTier(tierA)
        const hasTierB = isValidTier(tierB)

        // Items with tier come first, without tier come last
        if (hasTierA && !hasTierB) return -1
        if (!hasTierA && hasTierB) return 1

        // If both have tier, sort by tier priority (T0 > T1 > T2)
        if (hasTierA && hasTierB) {
          const priorityA = tierPriority[tierA] ?? 999
          const priorityB = tierPriority[tierB] ?? 999

          if (priorityA !== priorityB) {
            return priorityA - priorityB
          }
          // Same tier - sort by name (Chinese pinyin)
          return nameA.localeCompare(nameB, 'zh-CN')
        }

        // Neither has tier - sort by name only (Chinese pinyin)
        return nameA.localeCompare(nameB, 'zh-CN')
      })
  })

  return {
    codes,
    loading,
    searchQuery,
    selectedMode,
    selectedWeaponType,
    selectedDataSource,
    filteredCodes,
    groupedCodes,
    uniqueModes,
    loadCodes
  }
})
