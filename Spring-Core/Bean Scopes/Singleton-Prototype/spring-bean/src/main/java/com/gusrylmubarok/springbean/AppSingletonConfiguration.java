package com.gusrylmubarok.springbean;

import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

/**
 * @author Gusryl Mubarok
 * @email gusrylmubarok@gmail.com
 */

@Configuration
public class AppSingletonConfiguration {

    @Bean
    // @Scope("singleton")
    public UserService userService(){
        return new UserService();
    }

}
