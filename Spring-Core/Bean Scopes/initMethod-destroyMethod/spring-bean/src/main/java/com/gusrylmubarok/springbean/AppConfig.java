package com.gusrylmubarok.springbean;

import org.springframework.context.annotation.Bean;

/**
 * @author Gusryl Mubarok
 * @email gusrylmubarok@gmail.com
 */

public class AppConfig {
    @Bean(initMethod = "init", destroyMethod = "destroy")
    public DatabaseInitializr databaseInitiaizer() {
        return new DatabaseInitializr();
    }
}
