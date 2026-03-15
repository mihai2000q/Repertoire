import { ActionIcon, Group, Popover, Stack, Text, TextInput } from '@mantine/core'
import { IconPlus, IconSend } from '@tabler/icons-react'
import { useCreateSongArrangementMutation } from '../../../state/api/songsApi.ts'
import { toast } from 'react-toastify'
import { useEffect, useState } from 'react'

interface AddNewSongArrangementButtonProps {
  openedPopover: boolean
  setOpenedPopover: (value: boolean) => void
  songId: string
  onCreate?: (id: string) => void
}

function AddNewSongArrangementButton({
  openedPopover,
  setOpenedPopover,
  songId,
  onCreate
}: AddNewSongArrangementButtonProps) {
  const [createArrangement, { isLoading }] = useCreateSongArrangementMutation()

  const [error, setError] = useState(false)
  const [name, setName] = useState('')

  useEffect(() => setError(false), [openedPopover])

  function togglePopover() {
    setOpenedPopover(!openedPopover)
  }

  function handleChangeName(value: string) {
    setName(value)
    setError(value.trim() === '')
  }

  async function handleSubmit() {
    if (name.trim() === '') {
      setError(true)
      return
    }

    const res = await createArrangement({ name: name, songId: songId }).unwrap()

    toast.success('New song arrangement added!')
    onCreate?.(res.id)
    setOpenedPopover(false)
    setName('')
  }

  return (
    <Popover
      opened={openedPopover}
      onChange={setOpenedPopover}
      transitionProps={{ transition: 'fade-down' }}
      position={'bottom'}
      shadow={'md'}
      withArrow
      trapFocus
    >
      <Popover.Target>
        <ActionIcon
          variant={'grey'}
          aria-label={'add-new-arrangement'}
          size={'sm'}
          onClick={togglePopover}
        >
          <IconPlus size={14} />
        </ActionIcon>
      </Popover.Target>

      <Popover.Dropdown w={250}>
        <Stack gap={'xxs'}>
          <Text fz={'xs'} fw={500}>
            New Arrangement
          </Text>
          <Group gap={'xxs'}>
            <TextInput
              flex={1}
              size={'xs'}
              aria-label={'name'}
              placeholder={'Enter a name'}
              value={name}
              onChange={(e) => handleChangeName(e.currentTarget.value)}
              error={error}
            />
            <ActionIcon
              variant={'subtle'}
              aria-label={'submit'}
              size={'md'}
              loading={isLoading}
              disabled={error}
              onClick={handleSubmit}
            >
              <IconSend size={14} />
            </ActionIcon>
          </Group>
          {error && (
            <Text c={'red'} fz={'xxs'} fw={500} pl={'xxs'}>
              Name cannot be empty
            </Text>
          )}
        </Stack>
      </Popover.Dropdown>
    </Popover>
  )
}

export default AddNewSongArrangementButton
