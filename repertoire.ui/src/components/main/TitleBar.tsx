import { ActionIcon, Avatar, Group, Space, Stack, Text } from '@mantine/core'
import { IconArrowDownLeft, IconArrowsDiagonal, IconX } from '@tabler/icons-react'
import useTitleBarHeight from '../../hooks/useTitleBarHeight'
import logo from '../../assets/logo.png'
import {useAppSelector} from "../../state/store.ts";
import {useDocumentTitle} from "@mantine/hooks";

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

  const height = useTitleBarHeight()

  const title = useAppSelector((state) => state.global.documentTitle)
  useDocumentTitle(title)

  return (
    <Stack
      data-testid="title-bar"
      bg={'white'}
      pos={'fixed'}
      top={0}
      w={'100%'}
      gap={0}
      style={{ zIndex: 100, WebkitAppRegion: 'drag' }}
    >
      <Group gap={0} h={height} px={'xs'}>
        <Avatar src={logo} size={25} />
        <Text c={'gray.7'} fw={600} pl={6}>
          {title}
        </Text>

        <Space flex={1} />

        <Group gap={4} style={{ WebkitAppRegion: 'no-drag' }}>
          <ActionIcon
            aria-label={'minimize'}
            variant={'subtle'}
            size={'lg'}
            onClick={handleMinimize}
          >
            <IconArrowDownLeft />
          </ActionIcon>
          <ActionIcon
            aria-label={'maximize'}
            variant={'subtle'}
            size={'lg'}
            onClick={handleMaximize}
          >
            <IconArrowsDiagonal />
          </ActionIcon>
          <ActionIcon
            aria-label={'close'}
            variant={'subtle'}
            size={'lg'}
            onClick={handleClose}
            sx={(theme) => ({
              '&:hover': { backgroundColor: theme.colors.red[4], color: theme.white }
            })}
          >
            <IconX />
          </ActionIcon>
        </Group>
      </Group>
    </Stack>
  )
}

export default TitleBar
