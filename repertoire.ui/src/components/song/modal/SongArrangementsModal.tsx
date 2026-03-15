import {
  ActionIcon,
  Button,
  Center,
  Group,
  Modal,
  NumberInput,
  ScrollArea,
  Stack,
  Text,
  TextInput,
  Tooltip
} from '@mantine/core'
import { SongArrangement } from '../../../types/models/Song.ts'
import {
  useBulkUpdateSongArrangementsMutation,
  useDeleteSongArrangementMutation,
  useGetSongArrangementsQuery,
  useUpdateDefaultSongArrangementMutation
} from '../../../state/api/songsApi.ts'
import {
  IconArrowBackUp,
  IconMinus,
  IconPlus,
  IconStar,
  IconStarFilled,
  IconTrash
} from '@tabler/icons-react'
import { useEffect, useState } from 'react'
import SongArrangementsMenu from './SongArrangementsMenu.tsx'
import SongArrangementsModalLoader from './SongArrangementsModalLoader.tsx'
import { toast } from 'react-toastify'
import { UpdateSongArrangementRequest } from '../../../types/requests/SongRequests.ts'
import { useDisclosure, useElementSize, useMap } from '@mantine/hooks'
import AddNewSongArrangementButton from './AddNewSongArrangementButton.tsx'
import plural from '../../../utils/plural.ts'
import LoadingOverlayDebounced from '../../@ui/loader/LoadingOverlayDebounced.tsx'
import WarningModal from '../../@ui/modal/WarningModal.tsx'

interface SongArrangementsModalProps {
  opened: boolean
  onClose: () => void
  songId: string
  defaultId?: string
}

function SongArrangementsModal({ opened, onClose, songId, defaultId }: SongArrangementsModalProps) {
  const [deleteArrangement, { isLoading: isDeleteLoading }] = useDeleteSongArrangementMutation()
  const [updateDefaultArrangement, { isLoading: isUpdateDefaultLoading }] =
    useUpdateDefaultSongArrangementMutation()
  const [updateArrangements, { isLoading: isUpdateLoading }] =
    useBulkUpdateSongArrangementsMutation()

  const {
    data: arrangements,
    isLoading,
    isFetching
  } = useGetSongArrangementsQuery({ songId: songId })

  const [selectedArrangement, setSelectedArrangement] = useState<SongArrangement>(null)
  const [internalArrangements, setInternalArrangements] = useState<Map<string, SongArrangement>>(
    new Map<string, SongArrangement>()
  )

  const [openedDelete, { open: openDelete, close: closeDelete }] = useDisclosure(false)
  const [openedAddPopover, setOpenedAddPopover] = useState(false)
  const [justCreatedId, setJustCreatedId] = useState<string>(null)
  const onCreate = (id: string) => setJustCreatedId(id)

  const { ref, height } = useElementSize()

  const hasChanged = useMap<string, boolean>()
  function refreshHasChanged(internalArrangements: Map<string, SongArrangement>) {
    hasChanged.clear()
    for (const arrangement of arrangements) {
      if (
        arrangement.name !== internalArrangements.get(arrangement.id).name.trim() ||
        JSON.stringify(arrangement.sectionOccurrences) !==
          JSON.stringify(internalArrangements.get(arrangement.id).sectionOccurrences)
      ) {
        hasChanged.set(arrangement.id, true)
      } else {
        hasChanged.delete(arrangement.id)
      }
    }
  }

  const arrangementErrors = useMap<string, boolean>()
  function refreshHasError(internalArrangements: Map<string, SongArrangement>) {
    arrangementErrors.clear()
    for (const [_, arrangement] of internalArrangements) {
      if (arrangement.name.trim() === '') {
        arrangementErrors.set(arrangement.id, true)
      } else {
        arrangementErrors.delete(arrangement.id)
      }
    }
  }

  const isDefault = defaultId === selectedArrangement?.id

  useEffect(() => {
    if (!arrangements) return
    if (arrangements.length > 0) {
      // on first render and when the song changes with another (upon navigation)
      if (!selectedArrangement || selectedArrangement.songId !== songId) {
        setSelectedArrangement(arrangements[0])
      } else {
        // when it is updated, we want the latest
        const newArrangement = arrangements.find((a) => a.id === selectedArrangement.id)
        setSelectedArrangement(newArrangement)
      }
    }

    // on first render
    // when the song changes with another (upon navigation)
    // and when the sections change
    if (
      internalArrangements.size === 0 ||
      selectedArrangement?.songId !== songId ||
      selectedArrangement?.sectionOccurrences.length !==
        arrangements.find((a) => a.id === selectedArrangement?.id)?.sectionOccurrences.length
    ) {
      const newArrangements = new Map<string, SongArrangement>()
      arrangements.forEach((a) => newArrangements.set(a.id, a))
      setInternalArrangements(newArrangements)
    }

    // when a new arrangement has just been created
    if (justCreatedId) {
      const newArrangement = arrangements.find((a) => a.id === justCreatedId)
      const newArrangements = new Map<string, SongArrangement>([...internalArrangements])
      newArrangements.set(newArrangement.id, newArrangement)
      setInternalArrangements(newArrangements)
      setSelectedArrangement(newArrangement)
      setJustCreatedId(null)
    }
  }, [arrangements])

  // Handlers
  async function handleUpdateDefault() {
    await updateDefaultArrangement({
      id: isDefault ? null : selectedArrangement.id,
      songId: songId
    }).unwrap()
  }

  function handleChangeName(value: string) {
    const newArrangements = new Map<string, SongArrangement>([...internalArrangements])
    const newArrangement = {
      ...internalArrangements.get(selectedArrangement.id),
      name: value
    }
    newArrangements.set(selectedArrangement.id, newArrangement)
    setInternalArrangements(newArrangements)
    refreshHasChanged(newArrangements)
    refreshHasError(newArrangements)
  }

  function handleDecreaseOccurrence(sectionId: string) {
    const sectionOccurrences = internalArrangements
      .get(selectedArrangement.id)
      .sectionOccurrences.find((so) => so.section.id === sectionId)
    if (typeof sectionOccurrences.occurrences === 'number') {
      handleChangeOccurrence(sectionId, sectionOccurrences.occurrences - 1)
    }
  }

  function handleIncreaseOccurrence(sectionId: string) {
    const sectionOccurrences = internalArrangements
      .get(selectedArrangement.id)
      .sectionOccurrences.find((so) => so.section.id === sectionId)
    if (typeof sectionOccurrences.occurrences === 'number') {
      handleChangeOccurrence(sectionId, sectionOccurrences.occurrences + 1)
    }
  }

  function handleChangeOccurrence(sectionId: string, occurrence: number | string) {
    const arrangement = internalArrangements.get(selectedArrangement.id)
    const newArrangement: SongArrangement = {
      ...arrangement,
      sectionOccurrences: [
        ...arrangement.sectionOccurrences.map((so) =>
          so.section.id === sectionId
            ? {
                ...so,
                occurrences: occurrence
              }
            : so
        )
      ]
    }
    const newArrangements = new Map<string, SongArrangement>([...internalArrangements])
    newArrangements.set(selectedArrangement.id, newArrangement)
    setInternalArrangements(newArrangements)
    refreshHasChanged(newArrangements)
  }

  function handleOnBlurOccurrence(sectionId: string) {
    const sectionOccurrences = internalArrangements
      .get(selectedArrangement.id)
      .sectionOccurrences.find((so) => so.section.id === sectionId)
    if (sectionOccurrences.occurrences.toString().trim() === '') {
      handleChangeOccurrence(sectionId, 0)
    }
  }

  function handleReset() {
    const newArrangements = new Map<string, SongArrangement>([...internalArrangements])
    newArrangements.set(selectedArrangement.id, selectedArrangement)
    setInternalArrangements(newArrangements)
    refreshHasChanged(newArrangements)
    refreshHasError(newArrangements)
  }

  async function handleDelete() {
    await deleteArrangement({
      id: selectedArrangement.id,
      songId: songId
    }).unwrap()

    toast.success('Song arrangement deleted!')
    setSelectedArrangement(arrangements.find((a) => a.id !== selectedArrangement.id))
    internalArrangements.delete(selectedArrangement.id)
    hasChanged.delete(selectedArrangement.id)
    arrangementErrors.delete(selectedArrangement.id)
  }

  async function handleUpdate() {
    const updateArrangementRequests: UpdateSongArrangementRequest[] = Array.from(
      internalArrangements.entries()
    )
      .filter(([arrangementId]) => hasChanged.has(arrangementId))
      .map(([_, arrangement]) => ({
        id: arrangement.id,
        name: arrangement.name.trim(),
        occurrences: arrangement.sectionOccurrences.map((so) => ({
          sectionId: so.section.id,
          occurrences: so.toString().trim() === '' ? 0 : (so.occurrences as number)
        }))
      }))
    await updateArrangements({
      songId: songId,
      requests: updateArrangementRequests
    }).unwrap()

    toast.info(`Song arrangement${plural(updateArrangementRequests)} updated!`)
    hasChanged.clear()
    arrangementErrors.clear()
  }

  if (isLoading || !arrangements)
    return <SongArrangementsModalLoader opened={opened} onClose={onClose} />

  if (arrangements.length === 0) {
    return (
      <Modal.Root opened={opened} onClose={onClose} size={470}>
        <Modal.Overlay />
        <Modal.Content>
          <Modal.Header pb={'xs'} mih={0}>
            <Group gap={'xxs'}>
              <Modal.Title>Song Arrangements</Modal.Title>
              <AddNewSongArrangementButton
                openedPopover={openedAddPopover}
                setOpenedPopover={setOpenedAddPopover}
                songId={songId}
              />
            </Group>
          </Modal.Header>

          <Modal.Body px={0}>
            <Center h={100}>
              <Text>There are no arrangements yet. Try to add one</Text>
            </Center>
          </Modal.Body>
        </Modal.Content>
      </Modal.Root>
    )
  }

  return (
    <>
      <Modal.Root opened={opened} onClose={onClose}>
        <Modal.Overlay />
        <Modal.Content ref={ref}>
          <Modal.Header pb={'xs'} mih={0}>
            <Group gap={'xxs'}>
              <Modal.Title>Song Arrangements</Modal.Title>
              <AddNewSongArrangementButton
                openedPopover={openedAddPopover}
                setOpenedPopover={setOpenedAddPopover}
                songId={songId}
                onCreate={onCreate}
              />
            </Group>
            <Group gap={'xxs'} wrap={'nowrap'}>
              <ActionIcon
                aria-label={isDefault ? 'unset-default' : 'set-default'}
                variant={'subtle'}
                size={'sm'}
                loading={isUpdateDefaultLoading}
                onClick={handleUpdateDefault}
              >
                {isDefault ? <IconStarFilled size={14} /> : <IconStar size={14} />}
              </ActionIcon>
              <SongArrangementsMenu
                arrangements={arrangements}
                internalArrangements={internalArrangements}
                selectedArrangement={selectedArrangement}
                setSelectedArrangement={setSelectedArrangement}
                songId={songId}
                isFetching={isFetching}
                openAddPopover={() => setOpenedAddPopover(true)}
                defaultId={defaultId}
              />
              <Modal.CloseButton />
            </Group>
          </Modal.Header>
          <Modal.Body px={0}>
            <LoadingOverlayDebounced visible={isFetching} loaderProps={{ type: 'bars' }} />

            <Stack>
              {selectedArrangement && internalArrangements.get(selectedArrangement.id) && (
                <Stack px={'sm'}>
                  <Center w={'100%'} pos={'relative'}>
                    <TextInput
                      w={'65%'}
                      variant={'default'}
                      aria-label={'name'}
                      maxLength={30}
                      error={
                        arrangementErrors.get(selectedArrangement.id) && 'The name cannot be empty'
                      }
                      value={internalArrangements.get(selectedArrangement.id).name}
                      onChange={(e) => handleChangeName(e.target.value)}
                      styles={(theme) => ({
                        input: {
                          height: '30px',
                          minHeight: 0,
                          padding: '0 8px 0 8px',
                          textAlign: 'center',
                          fontSize: theme.fontSizes.lg,
                          fontWeight: 500,
                          borderColor: arrangementErrors.get(selectedArrangement.id)
                            ? theme.colors.red[4]
                            : 'transparent',
                          '&:hover': {
                            borderColor: arrangementErrors.get(selectedArrangement.id)
                              ? theme.colors.red[4]
                              : theme.colors.gray[4]
                          },
                          '&:focus': {
                            borderColor: arrangementErrors.get(selectedArrangement.id)
                              ? theme.colors.red[4]
                              : theme.colors.primary[3]
                          }
                        }
                      })}
                    />

                    <Group gap={'xxs'} pos={'absolute'} right={0} pr={'xs'}>
                      <Tooltip.Group openDelay={500} closeDelay={200}>
                        <Tooltip
                          label={
                            hasChanged.has(selectedArrangement.id)
                              ? 'Reset occurrences'
                              : 'Occurrences cannot be reset'
                          }
                        >
                          <ActionIcon
                            aria-label={'reset'}
                            variant={'grey'}
                            size={'md'}
                            disabled={!hasChanged.has(selectedArrangement.id)}
                            onClick={handleReset}
                          >
                            <IconArrowBackUp size={16} />
                          </ActionIcon>
                        </Tooltip>

                        <Tooltip label={'Delete arrangement'}>
                          <ActionIcon
                            aria-label={'delete'}
                            variant={'subtle-red'}
                            size={'md'}
                            onClick={openDelete}
                          >
                            <IconTrash size={16} />
                          </ActionIcon>
                        </Tooltip>
                      </Tooltip.Group>
                    </Group>
                  </Center>

                  <ScrollArea.Autosize mah={'50vh'} scrollbars={'y'} scrollbarSize={7}>
                    <Stack px={'sm'} gap={'xs'}>
                      {internalArrangements.get(selectedArrangement.id).sectionOccurrences
                        .length === 0 ? (
                        <Center h={100}>
                          <Text>There are no sections. Try to add one</Text>
                        </Center>
                      ) : (
                        internalArrangements
                          .get(selectedArrangement.id)
                          .sectionOccurrences.map((sectionOccurrence) => (
                            <Group
                              key={sectionOccurrence.section.id}
                              aria-label={`section-${sectionOccurrence.section.name}`}
                              gap={'xxs'}
                            >
                              <TextInput
                                size={'xs'}
                                w={75}
                                aria-label={'section-type'}
                                value={sectionOccurrence.section.songSectionType.name}
                                readOnly={true}
                                styles={{ input: { cursor: 'default' } }}
                                mr={4}
                              />
                              <TextInput
                                size={'xs'}
                                flex={1}
                                aria-label={'section-name'}
                                value={sectionOccurrence.section.name}
                                readOnly={true}
                                styles={{ input: { cursor: 'default' } }}
                              />
                              <Group gap={'xxs'}>
                                <ActionIcon
                                  aria-label={'decrease-section-occurrences'}
                                  size={'sm'}
                                  variant={'subtle'}
                                  disabled={sectionOccurrence.occurrences === 0}
                                  onClick={() =>
                                    handleDecreaseOccurrence(sectionOccurrence.section.id)
                                  }
                                >
                                  <IconMinus size={16} />
                                </ActionIcon>
                                <NumberInput
                                  w={40}
                                  size={'xs'}
                                  aria-label={'section-occurrences'}
                                  value={sectionOccurrence.occurrences}
                                  onChange={(o) =>
                                    handleChangeOccurrence(sectionOccurrence.section.id, o)
                                  }
                                  allowDecimal={false}
                                  allowNegative={false}
                                  hideControls
                                  styles={{
                                    input: { textAlign: 'center' }
                                  }}
                                  onBlur={() =>
                                    handleOnBlurOccurrence(sectionOccurrence.section.id)
                                  }
                                />
                                <ActionIcon
                                  aria-label={'increase-section-occurrences'}
                                  size={'sm'}
                                  variant={'subtle'}
                                  onClick={() =>
                                    handleIncreaseOccurrence(sectionOccurrence.section.id)
                                  }
                                >
                                  <IconPlus size={16} />
                                </ActionIcon>
                              </Group>
                            </Group>
                          ))
                      )}
                    </Stack>
                  </ScrollArea.Autosize>
                </Stack>
              )}

              <Tooltip
                disabled={hasChanged.size > 0 && arrangementErrors.size === 0}
                label={
                  hasChanged.size === 0 ? (
                    'You need to make a change before saving'
                  ) : (
                    <Text c={'white'} fz={'sm'}>
                      There are errors on {arrangementErrors.size} arrangement
                      {plural(arrangementErrors.size)}
                    </Text>
                  )
                }
                position="bottom"
              >
                <Button
                  mx={'md'}
                  data-disabled={hasChanged.size === 0 || arrangementErrors.size > 0}
                  loading={isUpdateLoading}
                  onClick={handleUpdate}
                >
                  Save Changes
                </Button>
              </Tooltip>
            </Stack>
          </Modal.Body>
        </Modal.Content>
      </Modal.Root>

      <WarningModal
        opened={openedDelete}
        onClose={closeDelete}
        title={'Delete Arrangement'}
        description={'Are you sure you want to delete this arrangement?'}
        onYes={handleDelete}
        isLoading={isDeleteLoading}
        size={'sm'}
        centered={false}
        styles={{ content: { marginTop: `${Math.max(87, height / 2 - 75)}px` } }}
      />
    </>
  )
}

export default SongArrangementsModal
