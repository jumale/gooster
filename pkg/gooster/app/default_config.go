package app

const defaultConfig = `
modules:
  - '#id': workdir
    col: 0
    row: 1
    width: 1
    height: 3
    focus_key: Ctrl-W
    extensions:
      - '#id': navigate
      - '#id': sort
  
  - '#id': output
    col: 1
    row: 1
    width: 1
    height: 1
    extensions: []
  
  - '#id': prompt
    col: 1
    row: 2
    width: 1
    height: 1
    focused: true
    focus_key: Ctrl-F
    extensions: []
  
  - '#id': status
    col: 0
    row: 0
    width: 2
    height: 1
    extensions:
      - '#id': workdir

  - '#id': complete
    col: 1
    row: 3
    width: 1
    height: 1
    extensions:
      - '#id': bash_completion
`
