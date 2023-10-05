- Open the spreadshet
- Go to flights and airlines tab and download as .csv format files
- Open cmd as administrator (Run as Administrator)
- If use container copy file into container
    ```bash
    docker cp ./Flight\ DataSet.csv container_name_or_id:/path/in/container
    ```
- Login postgres as developer type this code to copy dataset on the table airlines
    \copy table_name from '~address local storage.csv~' with (DELIMITER ',', FORMAT CSV, HEADER)
    example : \copy airlines from 'D:\Project\Bangkit\dataset\airlines.csv' with (DELIMITER ',', FORMAT CSV, HEADER)
- Do the same on the table flights