import {Divider, Grid, Group, Skeleton, Stack} from '@mantine/core'

function SongLoader() {
  return (
    <Stack data-testid={'song-loader'}>
      <Group align={'start'}>
        <Skeleton radius={'lg'} w={150} h={150} />
        <Stack gap={4} pt={'10px'}>
          <Skeleton w={65} h={15} />
          <Skeleton w={250} h={45} my={4} />
          <Group gap={4}>
            <Skeleton radius={'50%'} w={35} h={35} />
            <Skeleton w={80} h={15} />
            <Skeleton w={110} h={15} />
            <Skeleton w={45} h={15} />
          </Group>
        </Stack>
      </Group>

      <Divider />

      <Grid align="start" mb={'lg'}>
        <Grid.Col span={{ sm: 12, md: 4.5 }}>
          <Stack>
            <Skeleton w={'100%'} h={220} />
            <Skeleton w={'100%'} h={140} />
            <Skeleton w={'100%'} h={130} />
          </Stack>
        </Grid.Col>

        <Grid.Col span={{ sm: 12, md: 7.5 }}>
          <Stack>
            <Skeleton w={'100%'} h={100} />
            <Skeleton w={'100%'} h={240} />
          </Stack>
        </Grid.Col>
      </Grid>
    </Stack>
  )
}

export default SongLoader
