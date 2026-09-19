import { useMemo } from 'react'
import { SongPart, SongSection } from '../../../../../../../types/models/Song.ts'

export default function useSelectedSectionParts(
  selectedIds: string[],
  sections: SongSection[]
): [SongSection[], SongPart[]] {
  const selectedSections = useMemo(
    () => sections.filter((s) => selectedIds.includes(`section-${s.id}`)),
    [sections, selectedIds]
  )
  const selectedSectionParts = useMemo(
    () =>
      sections.flatMap((s) => s.parts.filter((p) => selectedIds.includes(`part-${p.id}:${s.id}`))),
    [sections, selectedIds]
  )
  return [selectedSections, selectedSectionParts]
}
