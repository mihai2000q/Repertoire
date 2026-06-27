import {
  ActionIcon,
  alpha,
  Center,
  Group,
  LoadingOverlay,
  Menu,
  ScrollArea,
  Stack,
  Text,
  Tooltip
} from '@mantine/core'
import { IconChecklist, IconStarFilled } from '@tabler/icons-react'
import { useDisclosure } from '@mantine/hooks'
import MenuItemConfirmation from '../../../../../../../../components/menu/item/MenuItemConfirmation.tsx'
import { useAddCustomSongRehearsalMutation } from '../../../../../../../../state/api/songsApi.ts'
import { useGetSongArrangementsQuery } from '../../../arrangements/state/api/songArrangementsApi.ts'
import { toast } from 'react-toastify'

interface CustomRehearsalButtonProps {
  songId: string
  sectionsCount: number
  defaultSongArrangementId?: string
}

function CustomRehearsalButton({
  songId,
  sectionsCount,
  defaultSongArrangementId
}: CustomRehearsalButtonProps) {
  const [opened, { toggle, close }] = useDisclosure(false)

  const [addCustomRehearsal, { isLoading: isCustomLoading }] = useAddCustomSongRehearsalMutation()

  const {
    data: arrangements,
    isFetching,
    isLoading
  } = useGetSongArrangementsQuery({ songId: songId })

  async function handleAddCustomRehearsal(arrangementId: string) {
    await addCustomRehearsal({ id: songId, arrangementId: arrangementId }).unwrap()
    toast.success(`Custom rehearsal added!`)
    close()
  }

  return (
    <Menu
      opened={opened}
      onChange={toggle}
      transitionProps={{ transition: 'fade-right' }}
      position={'right'}
      shadow={'sm'}
      withArrow
    >
      <Menu.Target>
        <Tooltip
          label={
            isLoading
              ? 'Loading arrangements...'
              : arrangements?.length === 0
                ? 'To add a custom rehearsal, you need arrangements'
                : sectionsCount === 0
                  ? 'To add a custom rehearsal, you need sections'
                  : 'Add Custom Rehearsal'
          }
          disabled={opened}
        >
          <ActionIcon
            aria-label={'add-custom-rehearsal'}
            disabled={sectionsCount === 0 || arrangements?.length === 0}
            loading={isLoading}
            variant={'grey'}
            size={'sm'}
            onClick={toggle}
          >
            <IconChecklist size={16} />
          </ActionIcon>
        </Tooltip>
      </Menu.Target>

      <Menu.Dropdown>
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
        </ScrollArea.Autosize>
      </Menu.Dropdown>
    </Menu>
  )
}

export default CustomRehearsalButton
