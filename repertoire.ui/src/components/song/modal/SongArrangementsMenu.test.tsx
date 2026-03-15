import { emptySongArrangement, reduxRender } from '../../../test-utils.tsx'
import SongArrangementsMenu from './SongArrangementsMenu.tsx'
import { SongArrangement } from '../../../types/models/Song.ts'
import { screen } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'

describe('Song Arrangements Menu', () => {
  const arrangements: SongArrangement[] = [
    {
      ...emptySongArrangement,
      id: '1',
      name: 'Perfect Rehearsal'
    },
    {
      ...emptySongArrangement,
      id: '2',
      name: 'Partial Rehearsal'
    }
  ]

  const selectedArrangement = arrangements[0]

  it('should render', async () => {
    const user = userEvent.setup()

    reduxRender(
      <SongArrangementsMenu
        arrangements={arrangements}
        internalArrangements={new Map<string, SongArrangement>()}
        selectedArrangement={selectedArrangement}
        setSelectedArrangement={vi.fn()}
        songId={''}
        openAddPopover={vi.fn()}
      />
    )

    expect(screen.getByRole('button', { name: selectedArrangement.name })).toBeInTheDocument()

    await user.click(screen.getByRole('button', { name: selectedArrangement.name }))
    expect(await screen.findByRole('menu', { name: selectedArrangement.name })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: selectedArrangement.name })).toHaveAttribute(
      'data-active',
      'true'
    )
    arrangements.forEach((arrangement) => {
      expect(screen.getByRole('menuitem', { name: arrangement.name })).toBeInTheDocument()
    })
    expect(screen.getByRole('menuitem', { name: /new arrangement/i })).toBeInTheDocument()
  })

  it('should display the internal arrangement name when available', async () => {
    const user = userEvent.setup()

    const internalArrangement = {
      ...selectedArrangement,
      name: 'Some New Name'
    }

    reduxRender(
      <SongArrangementsMenu
        arrangements={arrangements}
        internalArrangements={
          new Map<string, SongArrangement>([[selectedArrangement.id, internalArrangement]])
        }
        selectedArrangement={selectedArrangement}
        setSelectedArrangement={vi.fn()}
        songId={''}
        openAddPopover={vi.fn()}
      />
    )

    expect(screen.getByRole('button', { name: internalArrangement.name })).toBeInTheDocument()
    expect(screen.queryByRole('button', { name: selectedArrangement.name })).not.toBeInTheDocument()

    await user.click(screen.getByRole('button', { name: internalArrangement.name }))
    expect(await screen.findByRole('menu', { name: internalArrangement.name })).toBeInTheDocument()
    expect(screen.queryByRole('menu', { name: selectedArrangement.name })).not.toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: internalArrangement.name })).toBeInTheDocument()
    expect(
      screen.queryByRole('menuitem', { name: selectedArrangement.name })
    ).not.toBeInTheDocument()
  })

  it("should display invalid name when the internal arrangement's name is available, but invalid", async () => {
    const user = userEvent.setup()

    const internalArrangement = {
      ...selectedArrangement,
      name: ''
    }
    const invalidName = /invalid name/i

    reduxRender(
      <SongArrangementsMenu
        arrangements={arrangements}
        internalArrangements={
          new Map<string, SongArrangement>([[selectedArrangement.id, internalArrangement]])
        }
        selectedArrangement={selectedArrangement}
        setSelectedArrangement={vi.fn()}
        songId={''}
        openAddPopover={vi.fn()}
      />
    )

    expect(screen.getByRole('button', { name: invalidName })).toBeInTheDocument()

    await user.click(screen.getByRole('button', { name: invalidName }))
    expect(await screen.findByRole('menu', { name: invalidName })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: invalidName })).toBeInTheDocument()
  })

  it('should be able to switch selected arrangement when click menu items', async () => {
    const user = userEvent.setup()

    const newArrangement = arrangements[1]
    const setSelectedArrangement = vi.fn()

    const [{ rerender }] = reduxRender(
      <SongArrangementsMenu
        arrangements={arrangements}
        internalArrangements={new Map<string, SongArrangement>()}
        selectedArrangement={selectedArrangement}
        setSelectedArrangement={setSelectedArrangement}
        songId={''}
        openAddPopover={vi.fn()}
      />
    )

    expect(screen.getByRole('button', { name: selectedArrangement.name })).toBeInTheDocument()

    await user.click(screen.getByRole('button', { name: selectedArrangement.name }))
    await user.click(await screen.findByRole('menuitem', { name: newArrangement.name }))
    expect(setSelectedArrangement).toHaveBeenCalledExactlyOnceWith(newArrangement)

    rerender(
      <SongArrangementsMenu
        arrangements={arrangements}
        internalArrangements={new Map<string, SongArrangement>()}
        selectedArrangement={newArrangement}
        setSelectedArrangement={setSelectedArrangement}
        songId={''}
        openAddPopover={vi.fn()}
      />
    )

    expect(screen.getByRole('button', { name: newArrangement.name })).toBeInTheDocument()
  })

  it('should open add popover when clicking on new arrangement menu item', async () => {
    const user = userEvent.setup()

    const openAddPopover = vi.fn()

    reduxRender(
      <SongArrangementsMenu
        arrangements={arrangements}
        internalArrangements={new Map<string, SongArrangement>()}
        selectedArrangement={selectedArrangement}
        setSelectedArrangement={vi.fn()}
        songId={''}
        openAddPopover={openAddPopover}
      />
    )

    expect(screen.getByRole('button', { name: selectedArrangement.name })).toBeInTheDocument()

    await user.click(screen.getByRole('button', { name: selectedArrangement.name }))
    await user.click(await screen.findByRole('menuitem', { name: /new arrangement/i }))
    expect(openAddPopover).toHaveBeenCalledOnce()
  })

  it.skip('should be able to reorder arrangements', () => {})
})
