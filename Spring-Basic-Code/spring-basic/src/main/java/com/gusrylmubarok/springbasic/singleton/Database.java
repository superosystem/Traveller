package com.gusrylmubarok.springbasic.singleton;

/**
 * @author Gusryl Mubarok
 * @email gusrylmubarok@gmail.com
 */

public class Database {

    private static Database database;

    public static Database getInstance() {
        if (database == null) {
            database = new Database();
        }
        return database;
    }

    private Database() {

    }

}
