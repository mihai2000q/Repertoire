import { Center, Grid, Group, Skeleton, Stack } from '@mantine/core'

function HomeRecentlyPlayedLoader() {
  return (
    <>
      {Array.from(Array(20)).map((_, i) => (
        <Group key={i} pl={'lg'} pr={'xxs'} py={'xs'}>
          <Skeleton
            radius={'md'}
            h={38}
            w={38}
            style={(theme) => ({ boxShadow: theme.shadows.sm })}
          />
          <Grid flex={1} columns={12} align={'center'}>
            <Grid.Col span={{ base: 5, md: 8, xxl: 5 }}>
              <Stack gap={0}>
                <Skeleton w={125} h={15} mb={4} />
                <Skeleton w={70} h={10} />
              </Stack>
            </Grid.Col>
            <Grid.Col span={4} display={{ base: 'block', md: 'none', xxl: 'block' }}>
              <Skeleton flex={1} h={12} px={'xs'} />
            </Grid.Col>
            <Grid.Col span={{ base: 3, md: 4, xxl: 3 }}>
              <Center>
                <Skeleton w={50} h={15} px={'md'} />
              </Center>
            </Grid.Col>
          </Grid>
        </Group>
      ))}
    </>
  )
}

export default HomeRecentlyPlayedLoader
