import Song from '../../../types/models/Song.ts'
import {
  Button,
  Checkbox,
  ComboboxItem,
  Group,
  LoadingOverlay,
  Modal,
  NumberInput,
  Space,
  Stack,
  Tooltip
} from '@mantine/core'
import { useUpdateSongMutation } from '../../../state/api/songsApi.ts'
import Difficulty from '../../../types/enums/Difficulty.ts'
import { MouseEvent, useState } from 'react'
import { useDidUpdate, useInputState } from '@mantine/hooks'
import GuitarTuningSelect from '../../@ui/form/select/GuitarTuningSelect.tsx'
import DifficultySelect from '../../@ui/form/select/DifficultySelect.tsx'
import CustomIconMetronome from '../../@ui/icons/CustomIconMetronome.tsx'
import { toast } from 'react-toastify'

interface EditSongInformationModalProps {
  song: Song
  opened: boolean
  onClose: () => void
}

function EditSongInformationModal({ song, opened, onClose }: EditSongInformationModalProps) {
  const [updateSongMutation, { isLoading }] = useUpdateSongMutation()

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
      ? {
          value: song.difficulty,
          label: Difficulty[song.difficulty]
        }
      : null
  )
  const [bpm, setBpm] = useInputState<string | number>(song.bpm)
  const [isRecorded, setIsRecorded] = useInputState(song.isRecorded)

  const hasChanged =
    (typeof bpm === 'number' && bpm !== song.bpm) ||
    (difficulty?.value !== song.difficulty && (song.difficulty !== null || difficulty !== null)) ||
    guitarTuning?.value !== song.guitarTuning?.id ||
    isRecorded !== song.isRecorded

  useDidUpdate(() => {
    setGuitarTuning(
      song.guitarTuning
        ? {
            value: song.guitarTuning.id,
            label: song.guitarTuning.name
          }
        : null
    )
    setDifficulty(
      song.difficulty
        ? {
            value: song.difficulty,
            label: Difficulty[song.difficulty]
          }
        : null
    )
    setBpm(song.bpm)
    setIsRecorded(song.isRecorded)
  }, [song])

  async function updateSong(e: MouseEvent) {
    if (!hasChanged) {
      e.preventDefault()
      return
    }

    const parsedBpm = typeof bpm === 'string' ? null : bpm

    await updateSongMutation({
      ...song,
      albumId: song.album?.id,
      artistId: song.artist?.id,
      id: song.id,
      difficulty: difficulty ? (difficulty.value as Difficulty) : null,
      guitarTuningId: guitarTuning?.value,
      bpm: parsedBpm,
      isRecorded: isRecorded
    }).unwrap()

    onClose()
    toast.info('Song information updated!')
  }

  return (
    <Modal opened={opened} onClose={onClose} title={'Edit Song Information'}>
      <LoadingOverlay visible={isLoading} loaderProps={{ type: 'bars' }} />

      <Stack p={'xs'}>
        <Group>
          <GuitarTuningSelect option={guitarTuning} onChange={setGuitarTuning} />

          <DifficultySelect option={difficulty} onChange={setDifficulty} />
        </Group>

        <Group>
          <NumberInput
            min={1}
            allowNegative={false}
            allowDecimal={false}
            leftSection={<CustomIconMetronome size={20} />}
            label="Bpm"
            placeholder="Enter Bpm"
            value={bpm}
            onChange={setBpm}
          />

          <Space flex={0.4} />

          <Checkbox
            mt={'18px'}
            label={'Recorded'}
            checked={isRecorded}
            onChange={setIsRecorded}
            size={'md'}
          />
        </Group>

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
    </Modal>
  )
}

export default EditSongInformationModal
