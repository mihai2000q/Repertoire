import Song from '../../../types/models/Song.ts'
import {
  Button,
  ComboboxItem,
  LoadingOverlay,
  Modal,
  NumberInput,
  Stack,
  Tooltip
} from '@mantine/core'
import { useUpdateSongMutation } from '../../../state/songsApi.ts'
import Difficulty from '../../../utils/enums/Difficulty.ts'
import { MouseEvent, useState } from 'react'
import { useInputState } from '@mantine/hooks'
import GuitarTuningsSelect from '../../form/select/GuitarTuningsSelect.tsx'
import DifficultySelect from '../../form/select/DifficultySelect.tsx'
import { IconBmp } from '@tabler/icons-react'

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
          <GuitarTuningsSelect option={guitarTuning} onChange={setGuitarTuning} />

          <DifficultySelect option={difficulty} onChange={setDifficulty} />

          <NumberInput
            flex={1}
            min={1}
            leftSection={<IconBmp size={20} />}
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
