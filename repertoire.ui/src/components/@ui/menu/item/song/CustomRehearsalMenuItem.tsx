import { IconChecklist, IconStarFilled } from '@tabler/icons-react'
import { toast } from 'react-toastify'
import {
  useAddCustomSongRehearsalMutation,
  useGetSongArrangementsQuery
} from '../../../../../state/api/songsApi.ts'
import { useDisclosure } from '@mantine/hooks'
import { alpha, Center, Group, LoadingOverlay, Menu, ScrollArea, Stack, Text } from '@mantine/core'
import MenuItemConfirmation from '../MenuItemConfirmation.tsx'

interface CustomRehearsalMenuItemProps {
  id: string
  closeMenu: () => void
  defaultSongArrangementId?: string
}

function CustomRehearsalMenuItem({
  id,
  closeMenu,
  defaultSongArrangementId
}: CustomRehearsalMenuItemProps) {
  const [opened, { open, close }] = useDisclosure(false)

  const [addCustomRehearsal, { isLoading: isCustomLoading }] = useAddCustomSongRehearsalMutation()

  const { data: arrangements, isFetching } = useGetSongArrangementsQuery(
    { songId: id },
    { skip: !opened }
  )

  async function handleAddCustomRehearsal(arrangementId: string) {
    await addCustomRehearsal({ id: id, arrangementId: arrangementId }).unwrap()
    toast.success(`Custom rehearsal added!`)
    closeMenu()
  }

  return (
    <Menu.Sub onOpen={open} onClose={close} openDelay={100} closeDelay={250}>
      <Menu.Sub.Target>
        <Menu.Sub.Item leftSection={<IconChecklist size={14} />}>Custom Rehearsal</Menu.Sub.Item>
      </Menu.Sub.Target>

      <Menu.Sub.Dropdown>
        <ScrollArea.Autosize mah={'max(250px, 50vh)'} scrollbars={'y'} scrollbarSize={5}>
          <LoadingOverlay visible={isFetching} />

          <Stack gap={2}>
            {arrangements?.map((arrangement) => (
              <MenuItemConfirmation
                key={arrangement.id}
                isLoading={isCustomLoading}
                onConfirm={() => handleAddCustomRehearsal(arrangement.id)}
                sx={(theme) => ({
                  transition: '0.2s',
                  color: theme.colors.gray[7],
                  '.mantine-Menu-itemLabel': { display: 'flex', gap: 8 },
                  '.DefaultIcon': { color: theme.colors.gray[6] },
                  '&:hover': {
                    color: theme.colors.gray[8],
                    '.DefaultIcon': { color: alpha(theme.colors.gray[7], 0.7) },
                    backgroundColor: alpha(theme.colors.gray[1], 1)
                  }
                })}
              >
                <Group gap={'xxs'}>
                  <Text
                    lineClamp={1}
                    fz={'sm'}
                    c={'inherit'}
                    fw={defaultSongArrangementId === arrangement.id ? 500 : 400}
                    pt={defaultSongArrangementId === arrangement.id ? 1 : 0}
                  >
                    {arrangement.name}
                  </Text>
                  {defaultSongArrangementId === arrangement.id && (
                    <Center className={'DefaultIcon'}>
                      <IconStarFilled size={12} />
                    </Center>
                  )}
                </Group>
              </MenuItemConfirmation>
            ))}
          </Stack>

          {arrangements?.length === 0 && (
            <Text fz={'xs'} c={'dimmed'} px={'xs'}>
              No arrangements found
            </Text>
          )}
        </ScrollArea.Autosize>
      </Menu.Sub.Dropdown>
    </Menu.Sub>
  )
}

export default CustomRehearsalMenuItem
