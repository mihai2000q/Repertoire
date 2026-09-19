import { renderHook } from '@testing-library/react'
import { emptySongPart, emptySongSection } from '../../../../../../../test-utils.tsx'
import useSelectedSectionParts from './useSelectedSectionParts.ts'

describe('use selected section parts', () => {
  it('should return selected sections and matching parts when selected ids include section and part ids', () => {
    const sections = [
      {
        ...emptySongSection,
        id: '1',
        parts: [
          { ...emptySongPart, id: 'p1' },
          { ...emptySongPart, id: 'p2' }
        ]
      },
      { ...emptySongSection, id: '2', parts: [{ ...emptySongPart, id: 'p3' }] },
      { ...emptySongSection, id: '3', parts: [{ ...emptySongPart, id: 'p4' }] }
    ]

    const selectedIds = ['section-1', 'part-p1:1', 'part-p3:2', 'section-3']

    const { result } = renderHook(() => useSelectedSectionParts(selectedIds, sections))

    expect(result.current[0]).toStrictEqual([sections[0], sections[2]])
    expect(result.current[1]).toStrictEqual([sections[0].parts[0], sections[1].parts[0]])
  })
})
