package com.gusrylmubarok.springbasic;

import com.gusrylmubarok.springbasic.singleton.Database;
import org.junit.jupiter.api.Assertions;
import org.junit.jupiter.api.Test;

/**
 * @author Gusryl Mubarok
 * @email gusrylmubarok@gmail.com
 */

public class DatabaseTest {

    @Test
    void testSingleton() {
        var database1 = Database.getInstance();
        var database2 = Database.getInstance();

        Assertions.assertSame(database1, database2);
    }

}
