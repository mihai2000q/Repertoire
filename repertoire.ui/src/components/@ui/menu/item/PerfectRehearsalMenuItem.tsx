import { IconCheck, IconChecks, IconX } from '@tabler/icons-react'
import { ActionIcon, alpha, Box, Group, Menu, Transition } from '@mantine/core'
import { MouseEvent, useState } from 'react'
import { toast } from 'react-toastify'
import { useAddPerfectSongRehearsalMutation } from '../../../../state/api/songsApi.ts'

interface PerfectRehearsalMenuItemProps {
  songId: string
}

function PerfectRehearsalMenuItem({ songId }: PerfectRehearsalMenuItemProps) {
  const [addPerfectRehearsal, { isLoading }] = useAddPerfectSongRehearsalMutation()

  const [openedControls, setOpenedControls] = useState(false)

  function handleClick(e: MouseEvent) {
    e.stopPropagation()
    setOpenedControls(true)
  }

  function handleCancel(e: MouseEvent) {
    e.stopPropagation()
    setOpenedControls(false)
  }

  async function handleAddPerfectRehearsal(e: MouseEvent) {
    e.stopPropagation()
    await addPerfectRehearsal({ id: songId }).unwrap()
    toast.success(`Perfect rehearsal added!`)
    setOpenedControls(false)
  }

  return (
    <Menu.Item
      component={'div'}
      leftSection={<IconChecks size={14} />}
      onClick={handleClick}
      closeMenuOnClick={false}
      style={(theme) => ({
        ...(openedControls && {
          transition: '0.325s',
          cursor: 'default',
          color: theme.colors.gray[5],
          backgroundColor: 'transparent'
        })
      })}
    >
      <Box pos={'relative'}>
        Perfect Rehearsal
        <Transition
          mounted={openedControls}
          transition='fade-left'
          duration={325}
          timingFunction="ease"
        >
          {(styles) => (
            <Group
              gap={0}
              pos={'absolute'}
              pl={'md'}
              py={4}
              top={-9}
              right={-5}
              style={(theme) => ({
                ...styles,
                borderRadius: '16px',
                background: `linear-gradient(to right, transparent, ${theme.white} 27%)`
              })}
            >
              <ActionIcon
                aria-label={'cancel'}
                variant={'subtle'}
                size={26}
                radius={'50%'}
                disabled={isLoading}
                onClick={handleCancel}
                sx={(theme) => ({
                  color: theme.colors.red[4],
                  '&:hover': {
                    color: theme.colors.red[5],
                    backgroundColor: alpha(theme.colors.red[2], 0.35)
                  },
                  '&[data-disabled]': {
                    color: theme.colors.gray[4],
                    backgroundColor: 'transparent'
                  }
                })}
              >
                <IconX size={14} />
              </ActionIcon>
              <ActionIcon
                aria-label={'confirm'}
                variant={'subtle'}
                size={26}
                radius={'50%'}
                c={'green'}
                loading={isLoading}
                onClick={handleAddPerfectRehearsal}
                sx={(theme) => ({
                  '&:hover': { backgroundColor: alpha(theme.colors.green[2], 0.35) }
                })}
              >
                <IconCheck size={14} />
              </ActionIcon>
            </Group>
          )}
        </Transition>
      </Box>
    </Menu.Item>
  )
}

export default PerfectRehearsalMenuItem
