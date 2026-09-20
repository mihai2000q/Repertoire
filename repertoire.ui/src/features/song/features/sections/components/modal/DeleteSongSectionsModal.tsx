import WarningModal from '../../../../../../components/modal/WarningModal.tsx'
import { useBulkDeleteSongSectionsMutation } from '../../state/api/songSectionsApi.ts'
import { toast } from 'react-toastify'
import plural from '../../../../../../utils/plural.ts'
import { Checkbox, Stack, Text } from '@mantine/core'
import { useMemo, useState } from 'react'
import { SongSection } from '../../../../../../types/models/Song.ts'

interface DeleteSongSectionsModalProps {
  sections: SongSection[]
  songId: string
  opened: boolean
  onClose: () => void
  onDelete?: () => void
}

function DeleteSongSectionsModal({
  sections,
  songId,
  opened,
  onClose,
  onDelete
}: DeleteSongSectionsModalProps) {
  const [bulkDeleteSongSectionsMutation, { isLoading }] = useBulkDeleteSongSectionsMutation()
  const [deleteWithParts, setDeleteWithParts] = useState(false)

  const [ids, sectionPartIds] = useMemo(() => {
    const ids: string[] = []
    const sectionPartIds: Set<string> = new Set()
    sections.forEach((section) => {
      ids.push(section.id)
      section.parts.forEach((part) => {
        sectionPartIds.add(part.id)
      })
    })
    return [ids, Array.from(sectionPartIds)]
  }, [sections])

  async function handleDelete() {
    await bulkDeleteSongSectionsMutation({
      ids: ids,
      songId: songId,
      partIds: sectionPartIds
    }).unwrap()
    toast.success(
      `${ids.length} section${plural(ids.length)} deleted` +
        `${deleteWithParts ? ` with ${sectionPartIds.length} part${plural(sectionPartIds)}` : ''}!`
    )
    onDelete?.()
  }

  return (
    <WarningModal
      opened={opened}
      onClose={onClose}
      title={`Delete Section${plural(ids)}`}
      description={
        <Stack gap={5}>
          <Text>
            Are you sure you want to delete {ids.length} section{plural(ids)} ?
          </Text>
          <Checkbox
            checked={deleteWithParts}
            onChange={(event) => setDeleteWithParts(event.currentTarget.checked)}
            label={`Delete all associated parts (${sectionPartIds.length} part${plural(sectionPartIds)})`}
            c={'dimmed'}
            styles={{ label: { paddingLeft: 8 } }}
            disabled={sectionPartIds.length === 0}
          />
        </Stack>
      }
      onYes={handleDelete}
      isLoading={isLoading}
    />
  )
}

export default DeleteSongSectionsModal
