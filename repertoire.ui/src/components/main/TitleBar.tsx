import { ActionIcon, Avatar, Group, Space, Stack, Text } from '@mantine/core'
import { IconArrowDownLeft, IconArrowsDiagonal, IconX } from '@tabler/icons-react'

function TitleBar() {
  function handleMinimize() {
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-ignore
    window.electron.ipcRenderer.send('minimize')
  }

  function handleMaximize() {
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-ignore
    window.electron.ipcRenderer.send('maximize')
  }

  function handleClose() {
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-ignore
    window.electron.ipcRenderer.send('close')
  }

  return (
    <Stack
      data-testid="title-bar"
      bg={'white'}
      pos={'fixed'}
      top={0}
      w={'100%'}
      gap={0}
      style={{ WebkitAppRegion: 'drag' }}
    >
      <Group gap={0} h={45} px={'xs'} align={'center'}>
        <Avatar size={35} />
        <Text pl={4}>Repertoire</Text>

        <Space flex={1} />

        <Group gap={4} style={{ WebkitAppRegion: 'no-drag' }}>
          <ActionIcon aria-label={'minimize'} variant={'subtle'} size={'lg'} onClick={handleMinimize}>
            <IconArrowDownLeft />
          </ActionIcon>
          <ActionIcon aria-label={'maximize'} variant={'subtle'} size={'lg'} onClick={handleMaximize}>
            <IconArrowsDiagonal />
          </ActionIcon>
          <ActionIcon aria-label={'close'} variant={'subtle'} size={'lg'} onClick={handleClose}>
            <IconX />
          </ActionIcon>
        </Group>
      </Group>
    </Stack>
  )
}

export default TitleBar
