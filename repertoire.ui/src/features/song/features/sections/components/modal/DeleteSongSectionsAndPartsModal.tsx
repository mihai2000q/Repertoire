import { useBulkDeleteSongSectionsMutation } from '../../state/api/songSectionsApi.ts'
import { toast } from 'react-toastify'
import plural from '../../../../../../utils/plural.ts'
import { SongPart, SongSection } from '../../../../../../types/models/Song.ts'
import { Button, Center, Checkbox, Group, Modal, Stack, Text, Tooltip } from '@mantine/core'
import { useMemo } from 'react'
import SongSectionTypeBadge from '../section/SongSectionTypeBadge.tsx'
import RehearsalsBadge from '../../../outline/components/RehearsalsBadge.tsx'
import { IconInfoCircle } from '@tabler/icons-react'

interface DeleteSongSectionsAndPartsModalProps {
  sections: SongSection[]
  sectionParts: SongPart[]
  songId: string
  opened: boolean
  onClose: () => void
  onDelete?: () => void
}

function DeleteSongSectionsAndPartsModal({
  sections,
  sectionParts,
  songId,
  opened,
  onClose,
  onDelete
}: DeleteSongSectionsAndPartsModalProps) {
  const [bulkDeleteSongSectionsMutation, { isLoading }] = useBulkDeleteSongSectionsMutation()

  const [partIds, duplicatePartIds] = useMemo(() => {
    const newPartIds: Set<string> = new Set<string>()
    const newDuplicatePartIds: string[] = []
    sectionParts.forEach((p) => {
      if (newPartIds.has(p.id)) {
        newDuplicatePartIds.push(p.id)
      }
      newPartIds.add(p.id)
    })
    return [Array.from(newPartIds), newDuplicatePartIds]
  }, [sectionParts])

  async function handleDelete() {
    const ids = sections.map((s) => s.id)
    await bulkDeleteSongSectionsMutation({ ids: ids, songId: songId, partIds: partIds }).unwrap()
    toast.success(
      `${ids.length} section${plural(ids.length)} deleted with ${partIds.length} part${plural(partIds)}!`
    )
    onDelete?.()
    onClose()
  }

  return (
    <Modal
      opened={opened}
      onClose={onClose}
      title={`Delete Song Section${plural(sections)} with Part${plural(partIds)}`}
      centered
    >
      <Stack px={'xs'} py={0}>
        <Text fw={500}>Are you sure you want to delete the following:</Text>

        {sections.map((section) => (
          <Stack key={section.id}>
            <Checkbox.Card>
              <Group
                align={'start'}
                gap={'sm'}
                pl={'sm'}
                pr={'md'}
                sx={(theme) => ({ transition: '0.25s', border: `1px solid ${theme.colors.gray}` })}
              >
                <Checkbox.Indicator checked={true} />

                <Group gap={'xs'}>
                  <Text fw={600} fz={13} truncate={'end'}>
                    {section.name}
                  </Text>
                  <SongSectionTypeBadge songSectionType={section.songSectionType} />
                  <RehearsalsBadge rehearsals={section.rehearsals} />
                </Group>
              </Group>
            </Checkbox.Card>

            <Stack>
              {section.parts
                .filter((part) => partIds.some((id) => id === part.id))
                .map((part) => (
                  <Checkbox.Card key={part.id}>
                    <Group
                      align={'start'}
                      gap={'sm'}
                      pl={'45px'}
                      pr={'md'}
                      sx={(theme) => ({
                        transition: '0.25s',
                        border: `1px solid ${theme.colors.gray}`
                      })}
                    >
                      <Checkbox.Indicator checked={true} />

                      <Group gap={'xs'}>
                        <Text fw={600} fz={13} truncate={'end'}>
                          {part.name}
                        </Text>
                        <RehearsalsBadge rehearsals={part.rehearsals} />
                      </Group>

                      {duplicatePartIds.some((id) => id === part.id) && (
                        <Tooltip label={'This is a duplicate'}>
                          <Center>
                            <IconInfoCircle aria-label={`${part.name}-duplicate`} />
                          </Center>
                        </Tooltip>
                      )}
                    </Group>
                  </Checkbox.Card>
                ))}
            </Stack>
          </Stack>
        ))}

        <Group gap={'xxs'} style={{ alignSelf: 'end' }}>
          <Button variant={'subtle'} onClick={onClose}>
            Cancel
          </Button>
          <Button onClick={handleDelete} loading={isLoading}>
            Yes
          </Button>
        </Group>
      </Stack>
    </Modal>
  )
}

export default DeleteSongSectionsAndPartsModal
