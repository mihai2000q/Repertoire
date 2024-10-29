import { ActionIcon, Avatar, Group, Space, Stack, Text } from '@mantine/core'
import { IconArrowDownLeft, IconArrowsDiagonal, IconX } from '@tabler/icons-react'

function TitleBar() {
  function handleMinimize() {
    window.electron.ipcRenderer.send('minimize')
  }

  function handleMaximize() {
    window.electron.ipcRenderer.send('maximize')
  }

  function handleClose() {
    window.electron.ipcRenderer.send('close')
  }

  return (
    <Stack
      bg={'white'}
      pos={'fixed'}
      top={0}
      w={'100%'}
      gap={0}
      style={{ webkitAppRegion: 'drag' }}
    >
      <Group gap={0} h={45} px={'xs'} align={'center'}>
        <Avatar />
        <Text>Repertoire</Text>

        <Space flex={1} />

        <Group gap={4} style={{ webkitAppRegion: 'no-drag' }}>
          <ActionIcon variant={'subtle'} size={'lg'} onClick={handleMinimize}>
            <IconArrowDownLeft />
          </ActionIcon>
          <ActionIcon variant={'subtle'} size={'lg'} onClick={handleMaximize}>
            <IconArrowsDiagonal />
          </ActionIcon>
          <ActionIcon variant={'subtle'} size={'lg'} onClick={handleClose}>
            <IconX />
          </ActionIcon>
        </Group>
      </Group>
    </Stack>
  )
}

export default TitleBar
