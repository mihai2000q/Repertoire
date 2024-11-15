import Song from '../../../types/models/Song.ts'
import {
  Button,
  ComboboxItem,
  Group,
  Loader,
  LoadingOverlay,
  Modal,
  NumberInput,
  Select,
  Stack,
  Text,
  Tooltip
} from '@mantine/core'
import { useGetGuitarTuningsQuery, useUpdateSongMutation } from '../../../state/songsApi.ts'
import Difficulty from '../../../utils/enums/Difficulty.ts'
import { MouseEvent, useState } from 'react'
import { useInputState } from '@mantine/hooks'

interface EditSongInformationModalProps {
  song: Song
  opened: boolean
  onClose: () => void
}

function EditSongInformationModal({ song, opened, onClose }: EditSongInformationModalProps) {
  const [updateSongMutation, { isLoading }] = useUpdateSongMutation()

  const { data: guitarTuningsData, isLoading: isGuitarTuningsLoading } = useGetGuitarTuningsQuery()
  const guitarTunings = guitarTuningsData?.map((guitarTuning) => ({
    value: guitarTuning.id,
    label: guitarTuning.name
  }))

  const difficulties = Object.entries(Difficulty).map(([key, value]) => ({
    value: value,
    label: key
  }))

  const [guitarTuning, setGuitarTuning] = useState<ComboboxItem>(
    song.guitarTuning
      ? {
          value: song.guitarTuning.id,
          label: song.guitarTuning.name
        }
      : null
  )
  const [difficulty, setDifficulty] = useState<ComboboxItem>(
    song.difficulty
      ? difficulties.find((d) => d.value === song.difficulty)
      : null
  )
  const [bpm, setBpm] = useInputState<string | number>(song.bpm)

  const hasChanged =
    (typeof bpm === 'number' && bpm !== song.bpm) ||
    (difficulty?.value !== song.difficulty && (song.difficulty !== null || difficulty !== null)) ||
    guitarTuning?.value !== song.guitarTuning?.id

  async function updateSong(e: MouseEvent) {
    if (!hasChanged) {
      e.preventDefault()
      return
    }

    const parsedBpm = typeof bpm === 'string' ? null : bpm

    await updateSongMutation({
      ...song,
      id: song.id,
      difficulty: difficulty?.value as Difficulty,
      guitarTuningId: guitarTuning?.value,
      bpm: parsedBpm
    }).unwrap()

    onClose()
  }

  return (
    <Modal opened={opened} onClose={onClose} title={'Edit Song Information'}>
      <Modal.Body p={'xs'}>
        <LoadingOverlay
          visible={isLoading}
          zIndex={1000}
          overlayProps={{ radius: 'sm', blur: 2 }}
        />

        <Stack>
          {isGuitarTuningsLoading ? (
            <Group gap={'xs'} flex={1.25}>
              <Loader size={25} />
              <Text fz={'sm'} c={'dimmed'}>
                Loading Tunings...
              </Text>
            </Group>
          ) : (
            <Select
              flex={1.25}
              label={'Guitar Tuning'}
              placeholder={'Select Guitar Tuning'}
              data={guitarTunings}
              value={guitarTuning ? guitarTuning.value : null}
              onChange={(_, option) => setGuitarTuning(option)}
              maxDropdownHeight={150}
              clearable
            />
          )}

          <Select
            flex={1}
            label={'Difficulty'}
            placeholder={'Select Difficulty'}
            data={difficulties}
            value={difficulty ? difficulty.value : null}
            onChange={(_, option) => setDifficulty(option)}
            clearable
          />

          <NumberInput
            flex={1}
            min={1}
            label="Bpm"
            placeholder="Enter Bpm"
            value={bpm}
            onChange={setBpm}
          />

          <Tooltip
            disabled={hasChanged}
            label={'You need to make a change before saving'}
            position="bottom"
          >
            <Button data-disabled={!hasChanged} onClick={updateSong}>
              Save
            </Button>
          </Tooltip>
        </Stack>
      </Modal.Body>
    </Modal>
  )
}

export default EditSongInformationModal
