import { Skeleton, Stack } from '@mantine/core'

function HomeAlbumsLoader() {
  return (
    <>
      {Array.from(Array(20)).map((_, i) => (
        <Stack key={i} gap={'xs'} align={'center'}>
          <Skeleton
            radius={'lg'}
            h={150}
            w={150}
            style={(theme) => ({ boxShadow: theme.shadows.md })}
          />
          <Stack gap={0} align={'center'}>
            <Skeleton w={100} h={15} mb={4} />
            <Skeleton w={60} h={10} />
          </Stack>
        </Stack>
      ))}
    </>
  )
}

export default HomeAlbumsLoader
