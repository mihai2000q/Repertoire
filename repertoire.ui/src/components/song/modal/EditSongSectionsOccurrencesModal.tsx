import {
  ActionIcon,
  Button,
  Group,
  LoadingOverlay,
  Menu,
  Modal,
  NumberInput,
  ScrollArea,
  Stack,
  TextInput,
  Tooltip
} from '@mantine/core'
import { SongSection } from '../../../types/models/Song.ts'
import {
  useUpdateSongSectionsOccurrencesMutation,
  useUpdateSongSectionsPartialOccurrencesMutation
} from '../../../state/api/songsApi.ts'
import { useMap } from '@mantine/hooks'
import { useEffect, useState } from 'react'
import { toast } from 'react-toastify'
import { IconMinus, IconPlus } from '@tabler/icons-react'

enum RehearsalType {
  PerfectRehearsal = 'Perfect Rehearsal',
  PartialRehearsal = 'Partial Rehearsal'
}

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
  const [updateOccurrences, { isLoading: perfectIsLoading }] =
    useUpdateSongSectionsOccurrencesMutation()
  const [updatePartialOccurrences, { isLoading: partialIsLoading }] =
    useUpdateSongSectionsPartialOccurrencesMutation()
  const isLoading = perfectIsLoading || partialIsLoading

  const [rehearsalType, setRehearsalType] = useState(RehearsalType.PerfectRehearsal)

  const perfectOccurrences = useMap<string, string | number>([])
  const partialOccurrences = useMap<string, string | number>([])
  const occurrences =
    rehearsalType === RehearsalType.PerfectRehearsal ? perfectOccurrences : partialOccurrences
  useEffect(() => {
    occurrences.clear()
    partialOccurrences.clear()
    for (const section of sections) {
      perfectOccurrences.set(section.id, section.occurrences)
      partialOccurrences.set(section.id, section.partialOccurrences)
    }
  }, [sections])

  const hasChanged =
    JSON.stringify(Array.from(perfectOccurrences.entries()).map(([key, value]) => key + value)) !==
      JSON.stringify(sections.map((section) => section.id + section.occurrences)) ||
    JSON.stringify(Array.from(partialOccurrences.entries()).map(([key, value]) => key + value)) !==
      JSON.stringify(sections.map((section) => section.id + section.partialOccurrences))

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
    const sectionPerfectOccurrences = Array.from(perfectOccurrences.entries()).map(
      ([key, value]) => ({
        id: key,
        occurrences: typeof value === 'string' ? 0 : value
      })
    )
    const sectionPartialOccurrences = Array.from(partialOccurrences.entries()).map(
      ([key, value]) => ({
        id: key,
        partialOccurrences: typeof value === 'string' ? 0 : value
      })
    )

    // Optimizing is not a property, because this feature is temporary
    // It is planned that in the feature this will be replaced with something more advanced
    await updateOccurrences({
      songId: songId,
      sections: sectionPerfectOccurrences
    }).unwrap()
    await updatePartialOccurrences({
      songId: songId,
      sections: sectionPartialOccurrences
    }).unwrap()

    onClose()
    toast.info("Sections' occurrences updated!")
  }

  return (
    <Modal.Root opened={opened} onClose={onClose}>
      <Modal.Overlay />
      <Modal.Content>
        <Modal.Header>
          <Modal.Title>Edit Sections&#39; Occurrences</Modal.Title>
          <Group gap={'xxs'} wrap={'nowrap'}>
            <Menu position={'bottom-end'}>
              <Menu.Target>
                <Button
                  variant={'subtle'}
                  size={'compact-xs'}
                  styles={{ section: { marginLeft: 4 } }}
                >
                  {rehearsalType}
                </Button>
              </Menu.Target>
              <Menu.Dropdown>
                <Menu.Item
                  fz={'xs'}
                  onClick={() => setRehearsalType(RehearsalType.PerfectRehearsal)}
                >
                  {RehearsalType.PerfectRehearsal}
                </Menu.Item>
                <Menu.Item
                  fz={'xs'}
                  onClick={() => setRehearsalType(RehearsalType.PartialRehearsal)}
                >
                  {RehearsalType.PartialRehearsal}
                </Menu.Item>
              </Menu.Dropdown>
            </Menu>
            <Modal.CloseButton />
          </Group>
        </Modal.Header>
        <Modal.Body px={0}>
          <LoadingOverlay visible={isLoading} loaderProps={{ type: 'bars' }} />

          <Stack>
            <ScrollArea.Autosize mah={'50vh'} scrollbars={'y'} scrollbarSize={7}>
              <Stack px={'md'} gap={'xs'}>
                {sections.map((section) => (
                  <Group key={section.id} aria-label={`section-${section.name}`} gap={'xxs'}>
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
                    <Group gap={'xxs'}>
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
            </ScrollArea.Autosize>

            <Tooltip
              disabled={hasChanged}
              label={'You need to make a change before saving'}
              position="bottom"
            >
              <Button mx={'md'} data-disabled={!hasChanged} onClick={handleUpdateOccurrences}>
                Save Changes
              </Button>
            </Tooltip>
          </Stack>
        </Modal.Body>
      </Modal.Content>
    </Modal.Root>
  )
}

export default EditSongSectionsOccurrencesModal
