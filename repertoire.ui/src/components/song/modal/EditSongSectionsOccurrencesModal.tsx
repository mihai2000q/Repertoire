import {
  ActionIcon,
  Button,
  Group,
  LoadingOverlay,
  Modal,
  NumberInput,
  Stack,
  TextInput,
  Tooltip
} from '@mantine/core'
import { SongSection } from '../../../types/models/Song.ts'
import { useUpdateSongSectionsOccurrencesMutation } from '../../../state/api/songsApi.ts'
import { useMap } from '@mantine/hooks'
import { useEffect } from 'react'
import { toast } from 'react-toastify'
import { IconMinus, IconPlus } from '@tabler/icons-react'

interface EditSongSectionsOccurrencesModalProps {
  opened: boolean
  onClose: () => void
  sections: SongSection[]
  songId: string
}

function EditSongSectionsOccurrencesModal({
  opened,
  onClose,
  sections,
  songId
}: EditSongSectionsOccurrencesModalProps) {
  const [updateOccurrences, { isLoading }] = useUpdateSongSectionsOccurrencesMutation()

  const occurrences = useMap<string, string | number>([])
  useEffect(() => {
    occurrences.clear()
    sections.forEach((section) => occurrences.set(section.id, section.occurrences))
  }, [sections])

  const hasChanged =
    JSON.stringify(Array.from(occurrences.entries()).map(([key, value]) => key + value)) !==
    JSON.stringify(sections.map((section) => section.id + section.occurrences))

  function handleDecrease(sectionId: string) {
    const occ = occurrences.get(sectionId)
    if (typeof occ === 'string') return
    if (occ === 0) return
    occurrences.set(sectionId, occ - 1)
  }

  function handleIncrease(sectionId: string) {
    const occ = occurrences.get(sectionId)
    if (typeof occ === 'string') return
    occurrences.set(sectionId, occ + 1)
  }

  async function handleUpdateOccurrences() {
    const sectionOccurrences = Array.from(occurrences.entries()).map(([key, value]) => ({
      id: key,
      occurrences: typeof value === 'string' ? 0 : value
    }))

    await updateOccurrences({
      songId: songId,
      sections: sectionOccurrences
    }).unwrap()

    onClose()
    toast.info("Sections' occurrences updated!")
  }

  return (
    <Modal opened={opened} onClose={onClose} title={"Edit Sections' Occurrences"}>
      <Modal.Body p={0}>
        <LoadingOverlay visible={isLoading} loaderProps={{ type: 'bars' }} />

        <Stack>
          <Stack gap={'xs'}>
            {sections.map((section) => (
              <Group key={section.id} aria-label={`section-${section.name}`} gap={4}>
                <TextInput
                  size={'xs'}
                  w={75}
                  aria-label={'type'}
                  value={section.songSectionType.name}
                  readOnly={true}
                  mr={4}
                />
                <TextInput
                  size={'xs'}
                  flex={1}
                  aria-label={'name'}
                  value={section.name}
                  readOnly={true}
                />
                <Group gap={4}>
                  <ActionIcon
                    aria-label={'decrease-occurrences'}
                    size={'sm'}
                    variant={'subtle'}
                    disabled={occurrences.get(section.id) === 0}
                    onClick={() => handleDecrease(section.id)}
                  >
                    <IconMinus size={16} />
                  </ActionIcon>
                  <NumberInput
                    w={40}
                    size={'xs'}
                    aria-label={'occurrences'}
                    value={occurrences.get(section.id)}
                    onChange={(o) => occurrences.set(section.id, o)}
                    allowDecimal={false}
                    allowNegative={false}
                    hideControls
                    styles={{
                      input: { textAlign: 'center' }
                    }}
                    onBlur={() =>
                      occurrences.get(section.id).toString().trim() === '' &&
                      occurrences.set(section.id, 0)
                    }
                  />
                  <ActionIcon
                    aria-label={'increase-occurrences'}
                    size={'sm'}
                    variant={'subtle'}
                    onClick={() => handleIncrease(section.id)}
                  >
                    <IconPlus size={16} />
                  </ActionIcon>
                </Group>
              </Group>
            ))}
          </Stack>

          <Tooltip
            disabled={hasChanged}
            label={'You need to make a change before saving'}
            position="bottom"
          >
            <Button data-disabled={!hasChanged} onClick={handleUpdateOccurrences}>
              Save Changes
            </Button>
          </Tooltip>
        </Stack>
      </Modal.Body>
    </Modal>
  )
}

export default EditSongSectionsOccurrencesModal
