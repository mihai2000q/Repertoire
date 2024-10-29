import {mantineRender} from "../../test-utils.tsx";
import SongsLoader from "./SongsLoader.tsx";

describe('Songs Loader', () => {
  it('should render and display 20 card loaders', () => {
    const { container } = mantineRender(<SongsLoader />)

    expect(container.querySelectorAll('.mantine-Card-root')).toHaveLength(20)
  })
})
