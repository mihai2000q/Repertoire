import { ReactElement, useRef, useState } from 'react'
import useFixedDocumentTitle from '../hooks/useFixedDocumentTitle.ts'
import { ActionIcon, Button, Card, Group, SimpleGrid, Stack, Title } from '@mantine/core'
import HomeAlbums from '../components/home/HomeAlbums.tsx'
import { IconChevronLeft, IconChevronRight } from '@tabler/icons-react'

enum TopEntity {
  Songs,
  Albums,
  Artists
}

const TabButton = ({
  children,
  selected,
  onClick
}: {
  children: string
  selected: boolean
  onClick: () => void
}) => (
  <Button
    size={'compact-sm'}
    px={'xs'}
    sx={(theme) => ({
      color: theme.colors.gray[6],
      backgroundColor: 'transparent',

      '&:hover': {
        color: theme.colors.gray[9],
        backgroundColor: theme.colors.gray[1]
      },

      ...(selected && {
        color: theme.colors.primary[4],
        backgroundColor: 'transparent',
        '&:hover': {
          color: theme.colors.primary[4],
          backgroundColor: theme.colors.primary[0]
        }
      }),

      '&:active': {
        transform: 'unset'
      }
    })}
    style={{ transition: '0.25s' }}
    onClick={onClick}
  >
    {children}
  </Button>
)

function Home(): ReactElement {
  useFixedDocumentTitle('Home')

  const [entity, setEntity] = useState(TopEntity.Albums)

  const topRef = useRef<HTMLDivElement>(null)

  const handleTopNav = (direction: 'left' | 'right') => {
    if (!topRef.current) return

    console.log(topRef.current.scrollLeft)
    if (direction === 'left') {
      topRef.current.scrollTo({ left: topRef.current.scrollLeft - 200, behavior: 'smooth' })
    } else {
      topRef.current.scrollTo({ left: topRef.current.scrollLeft + 200, behavior: 'smooth' })
    }
  }

  return (
    <Stack mih={640} h={'100%'}>
      <Stack>
        <Title px={'xl'} order={2} fw={800} lh={1}>
          Welcome Back
        </Title>

        <Group px={'xl'} gap={0} justify={'space-between'}>
          <Group gap={0}>
            <TabButton
              selected={entity === TopEntity.Songs}
              onClick={() => setEntity(TopEntity.Songs)}
            >
              Songs
            </TabButton>
            <TabButton
              selected={entity === TopEntity.Albums}
              onClick={() => setEntity(TopEntity.Albums)}
            >
              Albums
            </TabButton>
            <TabButton
              selected={entity === TopEntity.Artists}
              onClick={() => setEntity(TopEntity.Artists)}
            >
              Artists
            </TabButton>
          </Group>

          <Group gap={4}>
            <ActionIcon
              aria-label={'back-button'}
              size={'lg'}
              variant={'grey'}
              radius={'50%'}
              onClick={() => handleTopNav('left')}
            >
              <IconChevronLeft size={20} />
            </ActionIcon>

            <ActionIcon
              aria-label={'forward-button'}
              size={'lg'}
              variant={'grey'}
              radius={'50%'}
              onClick={() => handleTopNav('right')}
            >
              <IconChevronRight size={20} />
            </ActionIcon>
          </Group>
        </Group>

        {entity === TopEntity.Albums && <HomeAlbums ref={topRef} />}
      </Stack>

      <SimpleGrid px={'xl'} cols={3} h={'100%'} mb={'md'}>
        <Card bg={'gray.3'} radius={'lg'} shadow={'lg'}>
          Card 1
        </Card>
        <Card bg={'gray.3'} radius={'lg'} shadow={'lg'}>
          Card 2
        </Card>
        <Card bg={'gray.3'} radius={'lg'} shadow={'lg'}>
          Card 3
        </Card>
      </SimpleGrid>
    </Stack>
  )
}

export default Home
