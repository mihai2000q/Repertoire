import { Card, Group, ScrollArea, Stack, Text } from '@mantine/core'
import { useEffect, useRef } from 'react'
import LoadingOverlayDebounced from '../../../../../../components/loader/LoadingOverlayDebounced.tsx'
import { useMain } from '../../../../../../context/MainContext.tsx'
import { ClickSelectProvider } from '../../../../../../context/ClickSelectContext.tsx'
import SongParts from '../parts/SongParts.tsx'
import NewHorizontalCard from '../../../../../../components/card/NewHorizontalCard.tsx'
import AddNewSongPart from './components/AddNewSongPart.tsx'
import AddNewSongSection from './components/AddNewSongSection.tsx'
import SongSections from '../sections/SongSections.tsx'
import OutlineView from './types/enums/OutlineView.ts'
import SongOutlineToolbar from './components/SongOutlineToolbar.tsx'
import { useGetSongSectionsQuery } from '../sections/state/api/songSectionsApi.ts'
import { useGetSongPartsQuery } from '../parts/state/api/songPartsApi.ts'
import { useDisclosure } from '@mantine/hooks'
import SongOutlineWidgetLoader from './components/loader/SongOutlineWidgetLoader.tsx'
import SongSectionsLoader from './components/loader/SongSectionsLoader.tsx'
import SongPartsLoader from './components/loader/SongPartsLoader.tsx'
import { useSongContext } from '../../context/SongContext.tsx'
import { useSongOutlineContext } from './context/SongOutlineContext.tsx'

interface SongOutlineWidgetProps {
  isSongFetching?: boolean
}

function SongOutlineWidget({ isSongFetching }: SongOutlineWidgetProps) {
  const { songId } = useSongContext()
  const { view, showDetails } = useSongOutlineContext()

  const {
    data: sections,
    isLoading: isSectionsLoading,
    isFetching: isSectionsFetching
  } = useGetSongSectionsQuery(
    {
      songId: songId
    },
    { skip: view !== OutlineView.Sections }
  )
  const {
    data: parts,
    isLoading: isPartsLoading,
    isFetching: isPartsFetching
  } = useGetSongPartsQuery({ songId: songId }, { skip: view !== OutlineView.Parts })

  const [openedAdd, { toggle: toggleAdd, close: closeAdd }] = useDisclosure(false)
  useEffect(() => closeAdd(), [songId, view])

  const ref = useRef<HTMLDivElement>(null)
  const scrollableRef = useRef<HTMLDivElement>(null)
  const { mainScroll } = useMain()

  function scrollAddIntoView() {
    scrollableRef.current?.scrollTo({ top: scrollableRef.current.scrollHeight, behavior: 'smooth' })
    mainScroll.ref.current?.scrollTo({
      top: mainScroll.ref.current.scrollHeight,
      behavior: 'smooth'
    })
  }

  function scrollCardIntoView() {
    ref.current?.scrollIntoView({ behavior: 'smooth' })
  }

  if (
    (view === OutlineView.Parts && (isPartsLoading || !parts) && !sections) ||
    (view === OutlineView.Sections && (isSectionsLoading || !sections) && !parts)
  ) {
    return <SongOutlineWidgetLoader outlineView={view} />
  }

  return (
    <ClickSelectProvider data={view === OutlineView.Sections ? sections : parts}>
      <Card ref={ref} variant={'widget'} aria-label={'outline-widget'} p={0}>
        <Stack gap={0}>
          <LoadingOverlayDebounced
            visible={isSongFetching}
            timeout={750}
            loaderProps={{ size: 'lg' }}
          />

          <Group px={'md'} pt={'md'} pb={'sm'} gap={'xxs'}>
            <Text fw={600} inline>
              Outline
            </Text>

            <SongOutlineToolbar
              toggleAdd={toggleAdd}
              parts={parts}
              sectionParts={sections?.flatMap((s) => s.parts)}
              scrollIntoView={scrollCardIntoView}
            />
          </Group>

          <ScrollArea.Autosize
            viewportRef={scrollableRef}
            scrollbars={'y'}
            scrollbarSize={7}
            mah={(showDetails ? 2 : 1) * 383.35}
            style={{ transition: 'max-height 0.25s' }}
          >
            <Stack gap={0}>
              {view === OutlineView.Sections &&
                (isSectionsLoading || !sections ? (
                  <SongSectionsLoader />
                ) : (
                  <SongSections
                    sections={sections}
                    isSongFetching={isSongFetching}
                    isSectionsFetching={isSectionsFetching}
                  />
                ))}
              {view === OutlineView.Parts &&
                (isPartsLoading || !parts ? (
                  <SongPartsLoader />
                ) : (
                  <SongParts
                    parts={parts}
                    isSongFetching={isSongFetching}
                    isPartsFetching={isPartsFetching}
                  />
                ))}

              {view === OutlineView.Sections && sections?.length === 0 && (
                <NewHorizontalCard ariaLabel={'add-new-song-section-card'} onClick={toggleAdd}>
                  Add New Song Section
                </NewHorizontalCard>
              )}
              {view === OutlineView.Parts && parts?.length === 0 && (
                <NewHorizontalCard ariaLabel={'add-new-song-part-card'} onClick={toggleAdd}>
                  Add New Song Part
                </NewHorizontalCard>
              )}

              {view === OutlineView.Sections && (
                <AddNewSongSection
                  songId={songId}
                  opened={openedAdd}
                  onClose={toggleAdd}
                  scrollIntoView={scrollAddIntoView}
                />
              )}
              {view === OutlineView.Parts && (
                <AddNewSongPart
                  opened={openedAdd}
                  onClose={toggleAdd}
                  scrollIntoView={scrollAddIntoView}
                />
              )}
            </Stack>
          </ScrollArea.Autosize>
        </Stack>
      </Card>
    </ClickSelectProvider>
  )
}

export default SongOutlineWidget
