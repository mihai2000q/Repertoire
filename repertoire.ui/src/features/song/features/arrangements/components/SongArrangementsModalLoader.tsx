import { Center, Group, Modal, Skeleton, Stack } from '@mantine/core'

interface SongArrangementsModalLoaderProps {
  opened: boolean
  onClose: () => void
}

function SongArrangementsModalLoader({ opened, onClose }: SongArrangementsModalLoaderProps) {
  return (
    <Modal.Root data-testid={'song-arrangements-loader'} opened={opened} onClose={onClose}>
      <Modal.Overlay />
      <Modal.Content>
        <Modal.Header pb={'xs'} mih={0}>
          <Group gap={'xxs'}>
            <Modal.Title>
              <Skeleton w={150} h={22} />
            </Modal.Title>
            <Skeleton w={20} h={20} />
          </Group>
          <Group gap={'xxs'} wrap={'nowrap'}>
            <Skeleton w={20} h={20} />
            <Skeleton w={100} h={20} />
            <Modal.CloseButton />
          </Group>
        </Modal.Header>

        <Modal.Body px={0}>
          <Stack px={'sm'}>
            <Center w={'100%'} pos={'relative'}>
              <Skeleton w={220} h={30} />

              <Group gap={'xxs'} pos={'absolute'} right={0} pr={'xs'}>
                <Skeleton w={25} h={25} />
                <Skeleton w={25} h={25} />
              </Group>
            </Center>

            <Stack px={'sm'} gap={'xs'}>
              {Array.from(Array(3)).map((_, i) => (
                <Group key={i} gap={'xs'}>
                  <Skeleton w={70} h={25} />
                  <Skeleton flex={1} h={25} />
                  <Skeleton w={50} h={25} />
                </Group>
              ))}
            </Stack>

            <Skeleton w={'100%'} h={30} />
          </Stack>
        </Modal.Body>
      </Modal.Content>
    </Modal.Root>
  )
}

export default SongArrangementsModalLoader
